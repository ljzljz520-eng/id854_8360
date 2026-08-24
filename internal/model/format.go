package model

import (
	"fmt"
	"sort"
	"strings"
)

func DisplayUserType(value UserType) string {
	switch value {
	case UserLeader:
		return "团长"
	case UserActor:
		return "演员"
	case UserStage:
		return "舞监"
	case UserTicket:
		return "票务"
	default:
		return "未知"
	}
}

func DisplayStatus(value string) string {
	switch value {
	case StatusDraft:
		return "草稿"
	case StatusReady:
		return "就绪"
	case StatusActive:
		return "进行中"
	case StatusComplete:
		return "完成"
	case StatusCanceled:
		return "已取消"
	case PerformancePlanned:
		return "计划"
	case PerformanceOpen:
		return "售票中"
	case PerformanceClosed:
		return "已闭场"
	default:
		return "未知"
	}
}

func CanonicalMenuSet(menus []string) string {
	clean := NormalizeMenus(menus)
	sort.Strings(clean)
	return strings.Join(clean, ",")
}

func IsTerminalStatus(status string) bool {
	return status == StatusComplete || status == StatusCanceled || status == PerformanceClosed
}

func CapacityLabel(capacity, sold int) string {
	available := capacity - sold
	if available < 0 {
		available = 0
	}
	return fmt.Sprintf("%d/%d", sold, capacity-available)
}

func RoleLabel(role Role) string {
	state := "停用"
	if role.Active {
		state = "启用"
	}
	return fmt.Sprintf("%s(%s,%s)", role.Name, DisplayUserType(role.UserType), state)
}

func RehearsalLabel(value Rehearsal) string {
	return fmt.Sprintf("%s %s %02d:00-%02d:00 [%s]", value.Production, value.Room, value.StartSlot, value.EndSlot, DisplayStatus(value.Status))
}

func PerformanceLabel(value Performance) string {
	return fmt.Sprintf("%s @ %s %s %s", value.Title, value.Venue, DisplayStatus(value.Status), CapacityLabel(value.Capacity, value.Sold))
}
