package roles

import (
	"errors"
	"strings"
	"theatrecontrol/internal/model"
)

func validateRole(role model.Role) error {
	if strings.TrimSpace(role.ID) == "" {
		return errors.New("role id is required")
	}
	if strings.TrimSpace(role.Name) == "" {
		return errors.New("role name is required")
	}
	if !model.ValidUserType(role.UserType) {
		return errors.New("unsupported user type")
	}
	if len(role.Name) > 80 {
		return errors.New("role name is too long")
	}
	return nil
}

func defaultActive(role model.Role) model.Role {
	if !role.Active {
		role.Active = true
	}
	role.Menus = sanitizeMenus(role.Menus)
	return role
}

func roleMatches(role model.Role, query string) bool {
	if query == "" {
		return true
	}
	q := strings.ToLower(query)
	return strings.Contains(strings.ToLower(role.Name), q) || strings.Contains(strings.ToLower(string(role.UserType)), q)
}
