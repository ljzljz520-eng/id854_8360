package report

import (
	"fmt"
	"sort"
	"strings"
	"theatrecontrol/internal/model"
)

type StatusCount struct {
	Status string
	Count  int
}

func CountRoleMenus(roles []model.Role) map[string]int {
	result := map[string]int{}
	for _, role := range roles {
		for _, menu := range role.Menus {
			result[menu]++
		}
	}
	return result
}

func SortStatusCounts(counts map[string]int) []StatusCount {
	result := make([]StatusCount, 0, len(counts))
	for status, count := range counts {
		result = append(result, StatusCount{Status: status, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count == result[j].Count {
			return result[i].Status < result[j].Status
		}
		return result[i].Count > result[j].Count
	})
	return result
}

func StatusLine(counts map[string]int) string {
	parts := make([]string, 0, len(counts))
	for _, value := range SortStatusCounts(counts) {
		parts = append(parts, fmt.Sprintf("%s:%d", value.Status, value.Count))
	}
	return strings.Join(parts, " ")
}

func (s *Service) RoleMenuUsage() (string, error) {
	roles, err := s.roles()
	if err != nil {
		return "", err
	}
	return StatusLine(CountRoleMenus(roles)), nil
}

func (s *Service) PerformanceStatuses() (string, error) {
	performances, err := s.performances()
	if err != nil {
		return "", err
	}
	counts := map[string]int{}
	for _, item := range performances {
		counts[item.Status]++
	}
	return StatusLine(counts), nil
}

func (s *Service) WorkflowSnapshot() (map[string]string, error) {
	summary, err := s.Summary()
	if err != nil {
		return nil, err
	}
	roles, err := s.RoleMenuUsage()
	if err != nil {
		return nil, err
	}
	performances, err := s.PerformanceStatuses()
	if err != nil {
		return nil, err
	}
	return map[string]string{"entities": fmt.Sprintf("%d", len(summary.Counts)), "role_menus": roles, "performance_statuses": performances}, nil
}
