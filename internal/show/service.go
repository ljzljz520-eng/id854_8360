package show

import (
	"errors"
	"fmt"
	"theatrecontrol/internal/audit"
	"theatrecontrol/internal/model"
	"theatrecontrol/internal/store"
)

type Service struct {
	db    *store.DB
	audit *audit.Service
}

func NewService(db *store.DB, logger *audit.Service) *Service { return &Service{db: db, audit: logger} }

func (s *Service) CreatePerformance(item model.Performance) (model.Performance, error) {
	if err := validatePerformance(item); err != nil {
		return model.Performance{}, err
	}
	item.Status = performanceState(item.Status)
	if err := s.db.Put("performance", item.ID, item); err != nil {
		return model.Performance{}, err
	}
	if s.audit != nil {
		_ = s.audit.Record("stage", "create_performance", "performance", item.ID, item.Title)
	}
	return item, nil
}

func (s *Service) Get(id string) (model.Performance, error) {
	var item model.Performance
	err := s.db.Get("performance", id, &item)
	return item, err
}

func (s *Service) SellTickets(id string, amount int) (model.Performance, error) {
	if amount <= 0 {
		return model.Performance{}, errors.New("ticket amount must be positive")
	}
	item, err := s.Get(id)
	if err != nil {
		return model.Performance{}, err
	}
	if item.Status != model.PerformanceOpen {
		return model.Performance{}, errors.New("performance is not open")
	}
	if seatsAvailable(item) < amount {
		return model.Performance{}, fmt.Errorf("only %d seats remain", seatsAvailable(item))
	}
	item.Sold += amount
	if err := s.db.Put("performance", id, item); err != nil {
		return model.Performance{}, err
	}
	if s.audit != nil {
		_ = s.audit.Record("ticketing", "sell_tickets", "performance", id, fmt.Sprintf("amount=%d", amount))
	}
	return item, nil
}

func (s *Service) Transition(id, status string) error {
	if status != model.PerformancePlanned && status != model.PerformanceOpen && status != model.PerformanceClosed {
		return errors.New("invalid performance status")
	}
	item, err := s.Get(id)
	if err != nil {
		return err
	}
	if item.Status == model.PerformanceClosed {
		return errors.New("closed performance cannot transition")
	}
	item.Status = status
	return s.db.Put("performance", id, item)
}

func (s *Service) AvailableSeats(id string) (int, error) {
	item, err := s.Get(id)
	if err != nil {
		return 0, err
	}
	return seatsAvailable(item), nil
}

func decode(data []byte, target any) error { return jsonDecode(data, target) }
