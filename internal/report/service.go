package report

import (
	"fmt"
	"sort"
	"strings"
	"theatrecontrol/internal/model"
	"theatrecontrol/internal/store"
)

type Summary struct {
	Counts           map[string]int `json:"counts"`
	ActiveRoles      int            `json:"active_roles"`
	OpenPerformances int            `json:"open_performances"`
	AuditHeadline    string         `json:"audit_headline"`
}

type Service struct{ db *store.DB }

func NewService(db *store.DB) *Service { return &Service{db: db} }

func (s *Service) Summary() (Summary, error) {
	counts, err := s.db.Snapshot()
	if err != nil {
		return Summary{}, err
	}
	roles, err := s.roles()
	if err != nil {
		return Summary{}, err
	}
	performances, err := s.performances()
	if err != nil {
		return Summary{}, err
	}
	active, open := 0, 0
	for _, role := range roles {
		if role.Active {
			active++
		}
	}
	for _, performance := range performances {
		if performance.Status == model.PerformanceOpen {
			open++
		}
	}
	entries, err := s.auditEntries()
	if err != nil {
		return Summary{}, err
	}
	return Summary{Counts: counts, ActiveRoles: active, OpenPerformances: open, AuditHeadline: fmt.Sprintf("%d audit events", len(entries))}, nil
}

func (s *Service) roles() ([]model.Role, error) {
	_, values, err := s.db.List("role")
	if err != nil {
		return nil, err
	}
	result := []model.Role{}
	for _, data := range values {
		var v model.Role
		if decode(data, &v) == nil {
			result = append(result, v)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}
func (s *Service) performances() ([]model.Performance, error) {
	_, values, err := s.db.List("performance")
	if err != nil {
		return nil, err
	}
	result := []model.Performance{}
	for _, data := range values {
		var v model.Performance
		if decode(data, &v) == nil {
			result = append(result, v)
		}
	}
	return result, nil
}
func (s *Service) auditEntries() ([]model.AuditEntry, error) {
	_, values, err := s.db.List("audit")
	if err != nil {
		return nil, err
	}
	result := []model.AuditEntry{}
	for _, data := range values {
		var v model.AuditEntry
		if decode(data, &v) == nil {
			result = append(result, v)
		}
	}
	return result, nil
}
func (s *Service) RenderText() (string, error) {
	summary, err := s.Summary()
	if err != nil {
		return "", err
	}
	keys := make([]string, 0, len(summary.Counts))
	for key := range summary.Counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	lines := []string{fmt.Sprintf("active roles: %d", summary.ActiveRoles), fmt.Sprintf("open performances: %d", summary.OpenPerformances), summary.AuditHeadline}
	for _, key := range keys {
		lines = append(lines, fmt.Sprintf("%s=%d", key, summary.Counts[key]))
	}
	return strings.Join(lines, "\n"), nil
}
func decode(data []byte, target any) error { return jsonDecode(data, target) }
