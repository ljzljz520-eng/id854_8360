package roles

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"theatrecontrol/internal/audit"
	"theatrecontrol/internal/model"
	"theatrecontrol/internal/store"
)

type BanService struct {
	db    *store.DB
	audit *audit.Service
}

func NewBanService(db *store.DB, logger *audit.Service) *BanService {
	return &BanService{db: db, audit: logger}
}

func validateBanRule(rule model.BanRule) error {
	if strings.TrimSpace(rule.ID) == "" {
		return errors.New("ban rule id is required")
	}
	if strings.TrimSpace(rule.RoleID) == "" {
		return errors.New("ban rule role is required")
	}
	if !allowedMenus[rule.Menu] {
		return fmt.Errorf("menu %q cannot be banned", rule.Menu)
	}
	if strings.TrimSpace(rule.Reason) == "" {
		return errors.New("ban reason is required")
	}
	return nil
}

func normalizeTags(tags []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.TrimSpace(strings.ToLower(tag))
		if tag == "" || seen[tag] {
			continue
		}
		seen[tag] = true
		result = append(result, tag)
	}
	sort.Strings(result)
	return result
}

func (s *BanService) Save(rule model.BanRule) (model.BanRule, error) {
	if err := validateBanRule(rule); err != nil {
		return model.BanRule{}, err
	}
	if !rule.Enabled {
		rule.Enabled = true
	}
	rule.AuditTags = normalizeTags(rule.AuditTags)
	if err := s.db.Put("ban_rule", rule.ID, rule); err != nil {
		return model.BanRule{}, err
	}
	if s.audit != nil {
		_ = s.audit.Record("leader", "save_ban_rule", "ban_rule", rule.ID, rule.Menu)
	}
	return rule, nil
}

func (s *BanService) Get(id string) (model.BanRule, error) {
	var rule model.BanRule
	err := s.db.Get("ban_rule", id, &rule)
	return rule, err
}

func (s *BanService) ListForRole(roleID string) ([]model.BanRule, error) {
	_, values, err := s.db.List("ban_rule")
	if err != nil {
		return nil, err
	}
	result := make([]model.BanRule, 0, len(values))
	for _, data := range values {
		var rule model.BanRule
		if decodeBan(data, &rule) == nil && rule.RoleID == roleID {
			result = append(result, rule)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}

func (s *BanService) IsBanned(roleID, menu string) (bool, error) {
	rules, err := s.ListForRole(roleID)
	if err != nil {
		return false, err
	}
	for _, rule := range rules {
		if rule.Enabled && rule.Menu == menu {
			return true, nil
		}
	}
	return false, nil
}

func (s *BanService) SetEnabled(id string, enabled bool) error {
	rule, err := s.Get(id)
	if err != nil {
		return err
	}
	rule.Enabled = enabled
	if err = s.db.Put("ban_rule", id, rule); err != nil {
		return err
	}
	if s.audit != nil {
		return s.audit.Record("leader", "toggle_ban_rule", "ban_rule", id, fmt.Sprintf("enabled=%t", enabled))
	}
	return nil
}

func (s *BanService) Delete(id string) error {
	if _, err := s.Get(id); err != nil {
		return err
	}
	if err := s.db.Delete("ban_rule", id); err != nil {
		return err
	}
	if s.audit != nil {
		return s.audit.Record("leader", "delete_ban_rule", "ban_rule", id, "removed")
	}
	return nil
}

func decodeBan(data []byte, target *model.BanRule) error { return decodeJSON(data, target) }
