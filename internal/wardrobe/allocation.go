package wardrobe

import (
	"errors"
	"theatrecontrol/internal/model"
)

func validateCostume(item model.Costume) error {
	if item.ID == "" || item.Name == "" {
		return errors.New("costume identity is incomplete")
	}
	if item.Quantity < 0 {
		return errors.New("costume quantity cannot be negative")
	}
	return nil
}

func validateAllocation(item model.CostumeAllocation) error {
	if item.ID == "" || item.CostumeID == "" || item.ActorID == "" {
		return errors.New("allocation identity is incomplete")
	}
	if item.Quantity <= 0 {
		return errors.New("allocation quantity must be positive")
	}
	return nil
}

func allocationState(value string) string {
	if value == "" {
		return model.StatusDraft
	}
	if value == model.StatusDraft || value == model.StatusReady || value == model.StatusActive || value == model.StatusComplete {
		return value
	}
	return model.StatusCanceled
}
