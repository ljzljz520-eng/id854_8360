package store

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"go.etcd.io/bbolt"
	"theatrecontrol/internal/model"
)

type DB struct{ db *bbolt.DB }

func Open(path string) (*DB, error) {
	if path == "" {
		return nil, errors.New("database path is required")
	}
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	result := &DB{db: db}
	if err := result.initialize(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return result, nil
}

func OpenEphemeral(path string) (*DB, func(), error) {
	db, err := Open(path)
	if err != nil {
		return nil, func() {}, err
	}
	cleanup := func() { _ = db.Close(); _ = os.Remove(path) }
	return db, cleanup, nil
}

func (s *DB) initialize() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, name := range bucketNames {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *DB) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *DB) Put(kind, id string, value any) error {
	if id == "" {
		return errors.New("record id is required")
	}
	data, err := encode(value)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucketFor(kind)).Put([]byte(id), data) })
}

func (s *DB) Get(kind, id string, target any) error {
	if id == "" {
		return errors.New("record id is required")
	}
	return s.db.View(func(tx *bbolt.Tx) error {
		value := tx.Bucket(bucketFor(kind)).Get([]byte(id))
		if value == nil {
			return fmt.Errorf("%s %q not found", kind, id)
		}
		return decode(cloneBytes(value), target)
	})
}

func (s *DB) Delete(kind, id string) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucketFor(kind)).Delete([]byte(id)) })
}

func (s *DB) List(kind string) ([][]byte, [][]byte, error) {
	keys, values := make([][]byte, 0), make([][]byte, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketFor(kind)).ForEach(func(k, v []byte) error {
			keys = append(keys, cloneBytes(k))
			values = append(values, cloneBytes(v))
			return nil
		})
	})
	if err != nil {
		return nil, nil, err
	}
	order := make([]int, len(keys))
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool { return string(keys[order[i]]) < string(keys[order[j]]) })
	sortedKeys, sortedValues := make([][]byte, len(keys)), make([][]byte, len(values))
	for i, index := range order {
		sortedKeys[i], sortedValues[i] = keys[index], values[index]
	}
	return sortedKeys, sortedValues, nil
}

func (s *DB) Count(kind string) (int, error) {
	count := 0
	err := s.db.View(func(tx *bbolt.Tx) error { count = tx.Bucket(bucketFor(kind)).Stats().KeyN; return nil })
	return count, err
}

func (s *DB) Snapshot() (map[string]int, error) {
	result := make(map[string]int)
	for _, kind := range []string{"role", "rehearsal", "costume", "allocation", "performance", "ban_rule", "audit", "assignment", "ticket_order"} {
		count, err := s.Count(kind)
		if err != nil {
			return nil, err
		}
		result[kind] = count
	}
	return result, nil
}

func (s *DB) PutRole(role model.Role) error { return s.Put("role", role.ID, role) }
