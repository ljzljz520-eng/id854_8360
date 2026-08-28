package rehearsal

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

func (s *Service) Schedule(item model.Rehearsal) (model.Rehearsal, error) {
	if item.ID == "" || item.Production == "" || item.Room == "" {
		return model.Rehearsal{}, errors.New("rehearsal identity is incomplete")
	}
	if err := validateWindow(item.StartSlot, item.EndSlot); err != nil {
		return model.Rehearsal{}, err
	}
	items, err := s.List(item.Room)
	if err != nil {
		return model.Rehearsal{}, err
	}
	for _, existing := range items {
		if existing.ID != item.ID && existing.Status != model.StatusCanceled && overlaps(existing, item) {
			return model.Rehearsal{}, fmt.Errorf("room %s is occupied at %s", item.Room, slotLabel(item.StartSlot, item.EndSlot))
		}
	}
	item.Status = rehearsalState(item.Status)
	if err := s.db.Put("rehearsal", item.ID, item); err != nil {
		return model.Rehearsal{}, err
	}
	if s.audit != nil {
		_ = s.audit.Record(item.Leader, "schedule_rehearsal", "rehearsal", item.ID, slotLabel(item.StartSlot, item.EndSlot))
	}
	return item, nil
}

func (s *Service) Get(id string) (model.Rehearsal, error) {
	var item model.Rehearsal
	err := s.db.Get("rehearsal", id, &item)
	return item, err
}

func (s *Service) List(room string) ([]model.Rehearsal, error) {
	_, values, err := s.db.List("rehearsal")
	if err != nil {
		return nil, err
	}
	result := []model.Rehearsal{}
	for _, data := range values {
		var item model.Rehearsal
		if err := decode(data, &item); err == nil && (room == "" || item.Room == room) {
			result = append(result, item)
		}
	}
	return result, nil
}

func (s *Service) Transition(id, status string) error {
	if !model.ValidStatus(status) {
		return errors.New("invalid rehearsal status")
	}
	item, err := s.Get(id)
	if err != nil {
		return err
	}
	if item.Status == model.StatusComplete || item.Status == model.StatusCanceled {
		return errors.New("finished rehearsal cannot transition")
	}
	item.Status = status
	return s.db.Put("rehearsal", id, item)
}

func decode(data []byte, target any) error { return jsonDecode(data, target) }
