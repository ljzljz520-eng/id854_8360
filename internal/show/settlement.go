package show

import (
	"errors"
	"fmt"
	"sort"
	"theatrecontrol/internal/audit"
	"theatrecontrol/internal/model"
	"theatrecontrol/internal/store"
)

type Settlement struct {
	PerformanceID    string `json:"performance_id"`
	GrossTickets     int    `json:"gross_tickets"`
	CancelledTickets int    `json:"cancelled_tickets"`
	NetTickets       int    `json:"net_tickets"`
	Status           string `json:"status"`
}

type SettlementService struct {
	db    *store.DB
	audit *audit.Service
}

func NewSettlementService(db *store.DB, logger *audit.Service) *SettlementService {
	return &SettlementService{db: db, audit: logger}
}

func (s *SettlementService) Calculate(performanceID string) (Settlement, error) {
	if performanceID == "" {
		return Settlement{}, errors.New("performance id is required")
	}
	orders, err := (&TicketService{db: s.db, audit: s.audit}).List(performanceID)
	if err != nil {
		return Settlement{}, err
	}
	performance, err := (&Service{db: s.db, audit: s.audit}).Get(performanceID)
	if err != nil {
		return Settlement{}, err
	}
	result := Settlement{PerformanceID: performanceID, Status: "open", NetTickets: performance.Sold}
	for _, order := range orders {
		if order.Status == "cancelled" {
			result.CancelledTickets += order.Quantity
		}
	}
	result.GrossTickets = result.NetTickets + result.CancelledTickets
	if performance.Status == model.PerformanceClosed {
		result.Status = "closed"
	}
	return result, nil
}

func (s *SettlementService) Close(performanceID string) (Settlement, error) {
	service := &Service{db: s.db, audit: s.audit}
	performance, err := service.Get(performanceID)
	if err != nil {
		return Settlement{}, err
	}
	if performance.Status == model.PerformanceClosed {
		return Settlement{}, errors.New("performance already closed")
	}
	performance.Status = model.PerformanceClosed
	if err = s.db.Put("performance", performance.ID, performance); err != nil {
		return Settlement{}, err
	}
	if s.audit != nil {
		_ = s.audit.Record("ticketing", "close_settlement", "performance", performanceID, "closed")
	}
	return s.Calculate(performanceID)
}

func FormatSettlement(value Settlement) string {
	return fmt.Sprintf("%s gross=%d cancelled=%d net=%d status=%s", value.PerformanceID, value.GrossTickets, value.CancelledTickets, value.NetTickets, value.Status)
}

func MergeSettlements(values []Settlement) Settlement {
	result := Settlement{Status: "open"}
	sort.Slice(values, func(i, j int) bool { return values[i].PerformanceID < values[j].PerformanceID })
	for _, value := range values {
		result.GrossTickets += value.GrossTickets
		result.CancelledTickets += value.CancelledTickets
		result.NetTickets += value.NetTickets
		if value.Status == "closed" {
			result.Status = "closed"
		}
	}
	return result
}
