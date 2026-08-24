package show

import (
	"errors"
	"theatrecontrol/internal/model"
)

func validatePerformance(item model.Performance) error {
	if item.ID == "" || item.Title == "" || item.Venue == "" {
		return errors.New("performance identity is incomplete")
	}
	if item.Capacity < 0 || item.Sold < 0 || item.Sold > item.Capacity {
		return errors.New("invalid ticket figures")
	}
	return nil
}

func performanceState(value string) string {
	if value == "" {
		return model.PerformancePlanned
	}
	if value == model.PerformancePlanned || value == model.PerformanceOpen || value == model.PerformanceClosed {
		return value
	}
	return model.PerformanceClosed
}

func seatsAvailable(item model.Performance) int {
	available := item.Capacity - item.Sold
	if available < 0 {
		return 0
	}
	return available
}
