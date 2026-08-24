package wardrobe

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

func (s *Service) RegisterCostume(item model.Costume) (model.Costume, error) {
	if err := validateCostume(item); err != nil {
		return model.Costume{}, err
	}
	if item.Status == "" {
		item.Status = model.StatusReady
	}
	if err := s.db.Put("costume", item.ID, item); err != nil {
		return model.Costume{}, err
	}
	if s.audit != nil {
		_ = s.audit.Record("wardrobe", "register_costume", "costume", item.ID, item.Name)
	}
	return item, nil
}

func (s *Service) Allocate(item model.CostumeAllocation) (model.CostumeAllocation, error) {
	if err := validateAllocation(item); err != nil {
		return model.CostumeAllocation{}, err
	}
	var costume model.Costume
	if err := s.db.Get("costume", item.CostumeID, &costume); err != nil {
		return model.CostumeAllocation{}, err
	}
	used, err := s.Used(item.CostumeID)
	if err != nil {
		return model.CostumeAllocation{}, err
	}
	if used+item.Quantity > costume.Quantity {
		return model.CostumeAllocation{}, fmt.Errorf("costume %s stock exceeded", costume.ID)
	}
	item.Status = allocationState(item.Status)
	if err := s.db.Put("allocation", item.ID, item); err != nil {
		return model.CostumeAllocation{}, err
	}
	if s.audit != nil {
		_ = s.audit.Record(item.ActorID, "allocate_costume", "allocation", item.ID, costume.Name)
	}
	return item, nil
}

func (s *Service) Used(costumeID string) (int, error) {
	_, values, err := s.db.List("allocation")
	if err != nil {
		return 0, err
	}
	used := 0
	for _, data := range values {
		var item model.CostumeAllocation
		if decode(data, &item) == nil && item.CostumeID == costumeID && item.Status != model.StatusCanceled {
			used += item.Quantity
		}
	}
	return used, nil
}

func (s *Service) GetAllocation(id string) (model.CostumeAllocation, error) {
	var item model.CostumeAllocation
	err := s.db.Get("allocation", id, &item)
	return item, err
}

func (s *Service) Release(id string) error {
	item, err := s.GetAllocation(id)
	if err != nil {
		return err
	}
	if item.Status == model.StatusCanceled {
		return errors.New("allocation already released")
	}
	item.Status = model.StatusCanceled
	return s.db.Put("allocation", id, item)
}

func decode(data []byte, target any) error { return jsonDecode(data, target) }
