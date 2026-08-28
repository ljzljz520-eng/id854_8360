package wardrobe

import (
	"path/filepath"
	"testing"
	"theatrecontrol/internal/model"
	"theatrecontrol/internal/store"
)

func TestCostumeAllocationStock(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "wardrobe.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	service := NewService(db, nil)
	if _, err = service.RegisterCostume(model.Costume{ID: "coat", Name: "外套", Quantity: 2}); err != nil {
		t.Fatal(err)
	}
	if _, err = service.Allocate(model.CostumeAllocation{ID: "first", CostumeID: "coat", ActorID: "a", Quantity: 2}); err != nil {
		t.Fatal(err)
	}
	if _, err = service.Allocate(model.CostumeAllocation{ID: "second", CostumeID: "coat", ActorID: "b", Quantity: 1}); err == nil {
		t.Fatal("expected stock error")
	}
	if err = service.Release("first"); err != nil {
		t.Fatal(err)
	}
	if used, err := service.Used("coat"); err != nil || used != 0 {
		t.Fatalf("release did not free stock: %d %v", used, err)
	}
	ledger := NewLedgerService(db, nil)
	if _, err = ledger.Adjust(LedgerEntry{ID: "restock", CostumeID: "coat", Delta: 1, Reason: "维修归还"}); err != nil {
		t.Fatal(err)
	}
	if items, err := ledger.InventoryByStatus(model.StatusReady); err != nil || len(items) != 1 {
		t.Fatalf("inventory lookup failed: %v %d", err, len(items))
	}
	if fit := CheckFit(model.Costume{ID: "coat", Size: "M"}, "S"); !fit.Approved {
		t.Fatalf("fit should be approved: %+v", fit)
	}
}
