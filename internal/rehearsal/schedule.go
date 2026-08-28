package rehearsal

import (
	"errors"
	"fmt"
	"theatrecontrol/internal/model"
)

func validateWindow(start, end int) error {
	if start < 0 || end < 0 {
		return errors.New("slots cannot be negative")
	}
	if start >= end {
		return errors.New("start slot must precede end slot")
	}
	if end > 24 {
		return errors.New("end slot exceeds day")
	}
	return nil
}

func overlaps(a, b model.Rehearsal) bool { return a.StartSlot < b.EndSlot && b.StartSlot < a.EndSlot }

func rehearsalState(status string) string {
	if status == "" {
		return model.StatusDraft
	}
	if status == model.StatusDraft || status == model.StatusReady || status == model.StatusActive {
		return status
	}
	return model.StatusCanceled
}

func slotLabel(start, end int) string { return fmt.Sprintf("%02d:00-%02d:00", start, end) }
