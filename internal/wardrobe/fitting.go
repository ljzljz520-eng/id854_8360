package wardrobe

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"theatrecontrol/internal/model"
)

var sizeOrder = map[string]int{"XS": 1, "S": 2, "M": 3, "L": 4, "XL": 5, "XXL": 6}

type FittingNote struct {
	ActorID       string
	CostumeID     string
	RequestedSize string
	Approved      bool
	Reason        string
}

func NormalizeSize(value string) string { return strings.ToUpper(strings.TrimSpace(value)) }

func SizeCompatible(costume, actor string) bool {
	left, okLeft := sizeOrder[NormalizeSize(costume)]
	right, okRight := sizeOrder[NormalizeSize(actor)]
	if !okLeft || !okRight {
		return false
	}
	difference := left - right
	return difference >= -1 && difference <= 1
}

func CheckFit(item model.Costume, actorSize string) FittingNote {
	normalized := NormalizeSize(actorSize)
	note := FittingNote{CostumeID: item.ID, RequestedSize: normalized}
	if normalized == "" {
		note.Reason = "actor size is missing"
		return note
	}
	if SizeCompatible(item.Size, normalized) {
		note.Approved = true
		note.Reason = "size is within one step"
	} else {
		note.Reason = fmt.Sprintf("costume size %s is not compatible", item.Size)
	}
	return note
}

func ValidateSize(value string) error {
	normalized := NormalizeSize(value)
	if _, ok := sizeOrder[normalized]; !ok {
		return errors.New("unsupported costume size")
	}
	return nil
}

func SortBySize(items []model.Costume) []model.Costume {
	result := append([]model.Costume(nil), items...)
	sort.Slice(result, func(i, j int) bool {
		left, right := sizeOrder[NormalizeSize(result[i].Size)], sizeOrder[NormalizeSize(result[j].Size)]
		if left == right {
			return result[i].ID < result[j].ID
		}
		return left < right
	})
	return result
}

func (s *Service) FittingOptions(actorSize string) ([]model.Costume, error) {
	items, err := s.listCostumes()
	if err != nil {
		return nil, err
	}
	result := make([]model.Costume, 0)
	for _, item := range items {
		if CheckFit(item, actorSize).Approved && item.Status != model.StatusCanceled {
			result = append(result, item)
		}
	}
	return SortBySize(result), nil
}

func (s *Service) listCostumes() ([]model.Costume, error) {
	_, values, err := s.db.List("costume")
	if err != nil {
		return nil, err
	}
	result := make([]model.Costume, 0, len(values))
	for _, data := range values {
		var item model.Costume
		if decode(data, &item) == nil {
			result = append(result, item)
		}
	}
	return result, nil
}
