package show

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"theatrecontrol/internal/audit"
	"theatrecontrol/internal/model"
	"theatrecontrol/internal/store"
)

type TicketService struct {
	db    *store.DB
	audit *audit.Service
}

func NewTicketService(db *store.DB, logger *audit.Service) *TicketService {
	return &TicketService{db: db, audit: logger}
}

func validateOrder(order model.TicketOrder) error {
	if order.ID == "" || order.PerformanceID == "" || strings.TrimSpace(order.Buyer) == "" {
		return errors.New("ticket order identity is incomplete")
	}
	if order.Quantity <= 0 {
		return errors.New("ticket order quantity must be positive")
	}
	return nil
}

func orderState(value string) string {
	if value == "" {
		return "reserved"
	}
	if value == "reserved" || value == "paid" || value == "cancelled" {
		return value
	}
	return "reserved"
}

func (s *TicketService) PlaceOrder(order model.TicketOrder) (model.TicketOrder, error) {
	if err := validateOrder(order); err != nil {
		return model.TicketOrder{}, err
	}
	performance, err := (&Service{db: s.db, audit: s.audit}).Get(order.PerformanceID)
	if err != nil {
		return model.TicketOrder{}, err
	}
	if performance.Status != model.PerformanceOpen {
		return model.TicketOrder{}, errors.New("performance is not accepting orders")
	}
	if seatsAvailable(performance) < order.Quantity {
		return model.TicketOrder{}, fmt.Errorf("not enough seats for order %s", order.ID)
	}
	performance.Sold += order.Quantity
	order.Status = orderState(order.Status)
	if err = s.db.Put("performance", performance.ID, performance); err != nil {
		return model.TicketOrder{}, err
	}
	if err = s.db.Put("ticket_order", order.ID, order); err != nil {
		return model.TicketOrder{}, err
	}
	if s.audit != nil {
		_ = s.audit.Record(order.Buyer, "place_ticket_order", "ticket_order", order.ID, fmt.Sprintf("quantity=%d", order.Quantity))
	}
	return order, nil
}

func (s *TicketService) Get(id string) (model.TicketOrder, error) {
	var order model.TicketOrder
	err := s.db.Get("ticket_order", id, &order)
	return order, err
}

func (s *TicketService) Cancel(id string) error {
	order, err := s.Get(id)
	if err != nil {
		return err
	}
	if order.Status == "cancelled" {
		return errors.New("ticket order already cancelled")
	}
	performance, err := (&Service{db: s.db, audit: s.audit}).Get(order.PerformanceID)
	if err != nil {
		return err
	}
	if performance.Sold < order.Quantity {
		return errors.New("performance sold count is inconsistent")
	}
	performance.Sold -= order.Quantity
	order.Status = "cancelled"
	if err = s.db.Put("performance", performance.ID, performance); err != nil {
		return err
	}
	if err = s.db.Put("ticket_order", order.ID, order); err != nil {
		return err
	}
	if s.audit != nil {
		return s.audit.Record(order.Buyer, "cancel_ticket_order", "ticket_order", id, "released")
	}
	return nil
}

func (s *TicketService) List(performanceID string) ([]model.TicketOrder, error) {
	_, values, err := s.db.List("ticket_order")
	if err != nil {
		return nil, err
	}
	result := make([]model.TicketOrder, 0, len(values))
	for _, data := range values {
		var order model.TicketOrder
		if decodeOrder(data, &order) == nil && (performanceID == "" || order.PerformanceID == performanceID) {
			result = append(result, order)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}

func (s *TicketService) QuantityByBuyer(buyer string) (int, error) {
	orders, err := s.List("")
	if err != nil {
		return 0, err
	}
	total := 0
	for _, order := range orders {
		if order.Buyer == buyer && order.Status != "cancelled" {
			total += order.Quantity
		}
	}
	return total, nil
}

func decodeOrder(data []byte, target *model.TicketOrder) error { return decodeOrderJSON(data, target) }
