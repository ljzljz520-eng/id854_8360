package store

import (
	"errors"
	"fmt"
	"go.etcd.io/bbolt"
	"sort"
	"strings"
)

type Query struct {
	Kind    string
	Prefix  string
	Limit   int
	Reverse bool
}

func (q Query) Validate() error {
	if q.Kind == "" {
		return errors.New("query kind is required")
	}
	if q.Limit < 0 {
		return errors.New("query limit cannot be negative")
	}
	return nil
}

func (s *DB) Query(q Query) ([]string, error) {
	if err := q.Validate(); err != nil {
		return nil, err
	}
	keys, err := s.Keys(q.Kind)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(keys))
	for _, key := range keys {
		if q.Prefix != "" && !strings.HasPrefix(key, q.Prefix) {
			continue
		}
		result = append(result, key)
	}
	if q.Reverse {
		sort.Sort(sort.Reverse(sort.StringSlice(result)))
	}
	if q.Limit > 0 && len(result) > q.Limit {
		result = result[:q.Limit]
	}
	return result, nil
}

func (s *DB) DeleteKind(kind string, ids []string) error {
	if kind == "" {
		return errors.New("delete kind is required")
	}
	if len(ids) == 0 {
		return errors.New("delete ids are required")
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(bucketFor(kind))
		for _, id := range ids {
			if id == "" {
				return errors.New("delete id is empty")
			}
			if err := bucket.Delete([]byte(id)); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *DB) Describe(kind string) (string, error) {
	count, err := s.Count(kind)
	if err != nil {
		return "", err
	}
	keys, err := s.Keys(kind)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s count=%d keys=%s", kind, count, strings.Join(keys, ",")), nil
}

func (s *DB) ExistsAll(kind string, ids []string) (bool, error) {
	if len(ids) == 0 {
		return false, nil
	}
	for _, id := range ids {
		found, err := s.Has(kind, id)
		if err != nil {
			return false, err
		}
		if !found {
			return false, nil
		}
	}
	return true, nil
}

func (s *DB) PrefixCount(kind, prefix string) (int, error) {
	keys, err := s.Query(Query{Kind: kind, Prefix: prefix})
	if err != nil {
		return 0, err
	}
	return len(keys), nil
}
