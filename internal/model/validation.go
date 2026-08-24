package model

import (
	"errors"
	"fmt"
	"strings"
)

func ValidateRoleShape(role Role) error {
	if err := ValidateEntityID(role.ID); err != nil {
		return fmt.Errorf("role: %w", err)
	}
	if strings.TrimSpace(role.Name) == "" {
		return errors.New("role name is empty")
	}
	if !ValidUserType(role.UserType) {
		return errors.New("role user type is invalid")
	}
	return nil
}

func ValidateRehearsalShape(value Rehearsal) error {
	if err := ValidateEntityID(value.ID); err != nil {
		return fmt.Errorf("rehearsal: %w", err)
	}
	if strings.TrimSpace(value.Production) == "" || strings.TrimSpace(value.Room) == "" {
		return errors.New("rehearsal location is incomplete")
	}
	if value.StartSlot < 0 || value.EndSlot > 24 || value.StartSlot >= value.EndSlot {
		return errors.New("rehearsal window is invalid")
	}
	return nil
}

func ValidatePerformanceShape(value Performance) error {
	if err := ValidateEntityID(value.ID); err != nil {
		return fmt.Errorf("performance: %w", err)
	}
	if strings.TrimSpace(value.Title) == "" || strings.TrimSpace(value.Venue) == "" {
		return errors.New("performance identity is incomplete")
	}
	if value.Capacity < 0 || value.Sold < 0 || value.Sold > value.Capacity {
		return errors.New("performance capacity is invalid")
	}
	return nil
}

func MenuSetDescription(menus []string) string {
	canonical := CanonicalMenuSet(menus)
	if canonical == "" {
		return "无可见权限"
	}
	return "可见菜单: " + canonical
}

func ValidateCostumeShape(value Costume) error {
	if err := ValidateEntityID(value.ID); err != nil {
		return fmt.Errorf("costume: %w", err)
	}
	if strings.TrimSpace(value.Name) == "" {
		return errors.New("costume name is empty")
	}
	if value.Quantity < 0 {
		return errors.New("costume quantity is invalid")
	}
	return nil
}

func ValidateTicketOrderShape(value TicketOrder) error {
	if err := ValidateEntityID(value.ID); err != nil {
		return fmt.Errorf("ticket order: %w", err)
	}
	if err := ValidateEntityID(value.PerformanceID); err != nil {
		return fmt.Errorf("ticket order performance: %w", err)
	}
	if strings.TrimSpace(value.Buyer) == "" || value.Quantity <= 0 {
		return errors.New("ticket order details are invalid")
	}
	return nil
}
