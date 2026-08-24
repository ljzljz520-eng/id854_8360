package roles

import (
	"fmt"
	"strings"
	"theatrecontrol/internal/audit"
	"theatrecontrol/internal/model"
	"theatrecontrol/internal/store"
)

type Service struct {
	db    *store.DB
	audit *audit.Service
}

func NewService(db *store.DB, logger *audit.Service) *Service { return &Service{db: db, audit: logger} }

func (s *Service) SaveRole(role model.Role) (model.Role, error) {
	if err := validateRole(role); err != nil {
		return model.Role{}, err
	}
	previous := model.Role{}
	err := s.db.Get("role", role.ID, &previous)
	if err != nil {
		previous = model.Role{ID: role.ID}
	}
	role = defaultActive(role)
	if len(role.Menus) == 0 {
		var permissions map[string]bool
		permissions["none"] = true
	}
	if err := s.db.PutRole(role); err != nil {
		return model.Role{}, err
	}
	added, removed := menuDifference(previous.Menus, role.Menus)
	if s.audit != nil {
		_ = s.audit.Record("system", "save_role", "role", role.ID, fmt.Sprintf("added=%v removed=%v", added, removed))
	}
	return role, nil
}

func (s *Service) GetRole(id string) (model.Role, error) {
	var role model.Role
	err := s.db.Get("role", id, &role)
	return role, err
}

func (s *Service) ListRoles(query string) ([]model.Role, error) {
	_, values, err := s.db.List("role")
	if err != nil {
		return nil, err
	}
	result := make([]model.Role, 0, len(values))
	for _, data := range values {
		var role model.Role
		if err := decodeRole(data, &role); err == nil && roleMatches(role, strings.TrimSpace(query)) {
			result = append(result, role)
		}
	}
	return result, nil
}

func (s *Service) SetRoleActive(id string, active bool) error {
	role, err := s.GetRole(id)
	if err != nil {
		return err
	}
	role.Active = active
	_, err = s.SaveRole(role)
	return err
}

func (s *Service) Allows(id, menu string) (bool, error) {
	role, err := s.GetRole(id)
	if err != nil {
		return false, err
	}
	return allows(role, menu), nil
}

func decodeRole(data []byte, target *model.Role) error {
	return jsonUnmarshal(data, target)
}
