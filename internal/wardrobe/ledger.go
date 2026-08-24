package wardrobe

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"theatrecontrol/internal/audit"
	"theatrecontrol/internal/model"
	"theatrecontrol/internal/store"
)

type LedgerEntry struct {
	ID        string `json:"id"`
	CostumeID string `json:"costume_id"`
	Delta     int    `json:"delta"`
	Reason    string `json:"reason"`
	Balance   int    `json:"balance"`
}
type LedgerService struct {
	db    *store.DB
	audit *audit.Service
}

func NewLedgerService(db *store.DB, logger *audit.Service) *LedgerService {
	return &LedgerService{db: db, audit: logger}
}

func validateLedger(entry LedgerEntry) error {
	if entry.ID == "" || entry.CostumeID == "" {
		return errors.New("ledger identity is incomplete")
	}
	if entry.Delta == 0 {
		return errors.New("ledger delta cannot be zero")
	}
	if strings.TrimSpace(entry.Reason) == "" {
		return errors.New("ledger reason is required")
	}
	return nil
}

func (s *LedgerService) Adjust(entry LedgerEntry) (LedgerEntry, error) {
	if err := validateLedger(entry); err != nil {
		return LedgerEntry{}, err
	}
	var costume model.Costume
	if err := s.db.Get("costume", entry.CostumeID, &costume); err != nil {
		return LedgerEntry{}, err
	}
	if costume.Quantity+entry.Delta < 0 {
		return LedgerEntry{}, fmt.Errorf("adjustment would make %s stock negative", costume.ID)
	}
	costume.Quantity += entry.Delta
	entry.Balance = costume.Quantity
	if err := s.db.Put("costume", costume.ID, costume); err != nil {
		return LedgerEntry{}, err
	}
	if err := s.db.Put("assignment", "ledger-"+entry.ID, entry); err != nil {
		return LedgerEntry{}, err
	}
	if s.audit != nil {
		_ = s.audit.Record("wardrobe", "adjust_costume_stock", "costume", costume.ID, entry.Reason)
	}
	return entry, nil
}

func (s *LedgerService) History(costumeID string) ([]LedgerEntry, error) {
	_, values, err := s.db.List("assignment")
	if err != nil {
		return nil, err
	}
	result := make([]LedgerEntry, 0, len(values))
	for _, data := range values {
		var entry LedgerEntry
		if decodeLedger(data, &entry) == nil && entry.CostumeID == costumeID {
			result = append(result, entry)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}

func (s *LedgerService) SetStatus(costumeID, status string) error {
	if status != model.StatusReady && status != model.StatusActive && status != model.StatusComplete && status != model.StatusCanceled {
		return errors.New("invalid costume status")
	}
	var costume model.Costume
	if err := s.db.Get("costume", costumeID, &costume); err != nil {
		return err
	}
	costume.Status = status
	if err := s.db.Put("costume", costumeID, costume); err != nil {
		return err
	}
	return nil
}

func (s *LedgerService) InventoryByStatus(status string) ([]model.Costume, error) {
	_, values, err := s.db.List("costume")
	if err != nil {
		return nil, err
	}
	result := make([]model.Costume, 0, len(values))
	for _, data := range values {
		var item model.Costume
		if decodeLedger(data, &item) == nil && (status == "" || item.Status == status) {
			result = append(result, item)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}

func decodeLedger(data []byte, target any) error { return decodeLedgerJSON(data, target) }
