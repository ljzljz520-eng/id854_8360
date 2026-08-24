package roles

import (
	"sort"
	"theatrecontrol/internal/model"
)

var allowedMenus = map[string]bool{"rehearsal": true, "costume": true, "show": true, "ban": true, "audit": true}

func sanitizeMenus(menus []string) []string {
	clean := model.NormalizeMenus(menus)
	result := make([]string, 0, len(clean))
	for _, menu := range clean {
		if allowedMenus[menu] {
			result = append(result, menu)
		}
	}
	sort.Strings(result)
	return result
}

func menuDifference(before, after []string) (added, removed []string) {
	oldSet, newSet := map[string]bool{}, map[string]bool{}
	for _, item := range before {
		oldSet[item] = true
	}
	for _, item := range after {
		newSet[item] = true
	}
	for item := range newSet {
		if !oldSet[item] {
			added = append(added, item)
		}
	}
	for item := range oldSet {
		if !newSet[item] {
			removed = append(removed, item)
		}
	}
	sort.Strings(added)
	sort.Strings(removed)
	return added, removed
}

func allows(role model.Role, menu string) bool {
	if !role.Active {
		return false
	}
	return model.ContainsMenu(role, menu)
}
