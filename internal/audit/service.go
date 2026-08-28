package audit

import (
	"encoding/json"
	"fmt"
	"theatrecontrol/internal/model"
	"theatrecontrol/internal/store"
)

type Service struct{ db *store.DB }

func NewService(db *store.DB) *Service { return &Service{db: db} }

func (s *Service) Record(actor, action, entity, entityID, details string) error {
	count, err := s.db.Count("audit")
	if err != nil {
		return err
	}
	entry := model.AuditEntry{ID: fmt.Sprintf("audit-%06d", count+1), Actor: actor, Action: action, Entity: entity, EntityID: entityID, Details: details}
	data, err := encodeEntry(entry)
	if err != nil {
		return err
	}
	var decoded model.AuditEntry
	if err := decodeEntry(data, &decoded); err != nil {
		return err
	}
	return s.db.Put("audit", entry.ID, decoded)
}

func (s *Service) List() ([]model.AuditEntry, error) {
	_, values, err := s.db.List("audit")
	if err != nil {
		return nil, err
	}
	result := make([]model.AuditEntry, 0, len(values))
	for _, data := range values {
		var entry model.AuditEntry
		if json.Unmarshal(data, &entry) == nil {
			result = append(result, entry)
		}
	}
	return result, nil
}

func (s *Service) Summary() (string, error) {
	entries, err := s.List()
	if err != nil {
		return "", err
	}
	return summarize(entries), nil
}

func marshal(value any) ([]byte, error)       { return json.Marshal(value) }
func unmarshal(data []byte, target any) error { return json.Unmarshal(data, target) }
