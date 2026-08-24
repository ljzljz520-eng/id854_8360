package store

import (
	"errors"
	"fmt"
	"go.etcd.io/bbolt"
	"sort"
)

type Operation struct {
	Kind   string
	ID     string
	Value  any
	Remove bool
}

func ValidateOperation(operation Operation) error {
	if operation.Kind == "" {
		return errors.New("operation kind is required")
	}
	if operation.ID == "" {
		return errors.New("operation id is required")
	}
	if operation.Remove {
		return nil
	}
	if operation.Value == nil {
		return errors.New("operation value is required")
	}
	return nil
}

func (s *DB) ApplyBatch(operations []Operation) error {
	if len(operations) == 0 {
		return errors.New("batch cannot be empty")
	}
	encoded := make([][]byte, len(operations))
	for i, operation := range operations {
		if err := ValidateOperation(operation); err != nil {
			return fmt.Errorf("operation %d: %w", i, err)
		}
		if !operation.Remove {
			data, err := encode(operation.Value)
			if err != nil {
				return err
			}
			encoded[i] = data
		}
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		for i, operation := range operations {
			bucket := tx.Bucket(bucketFor(operation.Kind))
			if operation.Remove {
				if err := bucket.Delete([]byte(operation.ID)); err != nil {
					return err
				}
				continue
			}
			if err := bucket.Put([]byte(operation.ID), encoded[i]); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *DB) Has(kind, id string) (bool, error) {
	found := false
	err := s.db.View(func(tx *bbolt.Tx) error { found = tx.Bucket(bucketFor(kind)).Get([]byte(id)) != nil; return nil })
	return found, err
}

func (s *DB) Keys(kind string) ([]string, error) {
	keys, _, err := s.List(kind)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(keys))
	for _, key := range keys {
		result = append(result, string(key))
	}
	sort.Strings(result)
	return result, nil
}

func (s *DB) Copy(kind, source, target string) error {
	if source == target {
		return errors.New("source and target must differ")
	}
	var data []byte
	if err := s.GetRaw(kind, source, &data); err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucketFor(kind)).Put([]byte(target), data) })
}

func (s *DB) GetRaw(kind, id string, target *[]byte) error {
	return s.db.View(func(tx *bbolt.Tx) error {
		value := tx.Bucket(bucketFor(kind)).Get([]byte(id))
		if value == nil {
			return fmt.Errorf("%s %q not found", kind, id)
		}
		*target = cloneBytes(value)
		return nil
	})
}
