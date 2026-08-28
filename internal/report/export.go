package report

import (
	"fmt"
	"sort"
	"strings"
	"theatrecontrol/internal/model"
)

func roleMatrix(roles []model.Role) []string {
	lines := make([]string, 0, len(roles))
	sort.Slice(roles, func(i, j int) bool { return roles[i].ID < roles[j].ID })
	for _, role := range roles {
		menus := append([]string(nil), role.Menus...)
		sort.Strings(menus)
		state := "inactive"
		if role.Active {
			state = "active"
		}
		lines = append(lines, fmt.Sprintf("%s|%s|%s|%s", role.ID, role.Name, state, strings.Join(menus, ",")))
	}
	return lines
}

func (s *Service) RoleMatrix() (string, error) {
	roles, err := s.roles()
	if err != nil {
		return "", err
	}
	return strings.Join(roleMatrix(roles), "\n"), nil
}

func (s *Service) CapacityBoard() (string, error) {
	performances, err := s.performances()
	if err != nil {
		return "", err
	}
	lines := make([]string, 0, len(performances))
	sort.Slice(performances, func(i, j int) bool { return performances[i].ID < performances[j].ID })
	for _, item := range performances {
		lines = append(lines, fmt.Sprintf("%s|%s|%d/%d|%s", item.ID, item.Title, item.Sold, item.Capacity, item.Status))
	}
	return strings.Join(lines, "\n"), nil
}

func (s *Service) EntityCount(kind string) (int, error) { return s.db.Count(kind) }
