package model

const (
	StatusDraft        = "draft"
	StatusReady        = "ready"
	StatusActive       = "active"
	StatusComplete     = "complete"
	StatusCanceled     = "canceled"
	PerformancePlanned = "planned"
	PerformanceOpen    = "open"
	PerformanceClosed  = "closed"
)

func ValidUserType(value UserType) bool {
	switch value {
	case UserLeader, UserActor, UserStage, UserTicket:
		return true
	default:
		return false
	}
}

func ValidStatus(value string) bool {
	if value == StatusDraft || value == StatusReady || value == StatusActive {
		return true
	}
	return value == StatusComplete || value == StatusCanceled
}

func NormalizeMenus(menus []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(menus))
	for _, menu := range menus {
		if menu == "" || seen[menu] {
			continue
		}
		seen[menu] = true
		result = append(result, menu)
	}
	return result
}

func ContainsMenu(role Role, menu string) bool {
	for _, candidate := range role.Menus {
		if candidate == menu {
			return true
		}
	}
	return false
}
