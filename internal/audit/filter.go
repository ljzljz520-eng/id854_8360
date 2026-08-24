package audit

import (
	"sort"
	"strings"
	"theatrecontrol/internal/model"
)

func (s *Service) Query(entity, actor string) ([]model.AuditEntry, error) {
	entries, err := s.List()
	if err != nil {
		return nil, err
	}
	result := make([]model.AuditEntry, 0, len(entries))
	for _, entry := range entries {
		if entity != "" && entry.Entity != entity {
			continue
		}
		if actor != "" && entry.Actor != actor {
			continue
		}
		result = append(result, entry)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}

func (s *Service) ActionsFor(entityID string) ([]string, error) {
	entries, err := s.Query("", "")
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	for _, entry := range entries {
		if entry.EntityID == entityID {
			result = append(result, entry.Action)
		}
	}
	return result, nil
}

func (s *Service) Search(text string) ([]model.AuditEntry, error) {
	entries, err := s.List()
	if err != nil {
		return nil, err
	}
	query := strings.ToLower(strings.TrimSpace(text))
	result := make([]model.AuditEntry, 0)
	for _, entry := range entries {
		value := strings.ToLower(entry.Action + " " + entry.Entity + " " + entry.EntityID + " " + entry.Details)
		if query == "" || strings.Contains(value, query) {
			result = append(result, entry)
		}
	}
	return result, nil
}
