package show

import (
	"path/filepath"
	"testing"
	"theatrecontrol/internal/model"
	"theatrecontrol/internal/store"
)

func TestPerformanceTicketFlow(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "show.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	service := NewService(db, nil)
	if _, err = service.CreatePerformance(model.Performance{ID: "matinee", Title: "午后", Venue: "剧场", Capacity: 4}); err != nil {
		t.Fatal(err)
	}
	if err = service.Transition("matinee", model.PerformanceOpen); err != nil {
		t.Fatal(err)
	}
	if _, err = service.SellTickets("matinee", 3); err != nil {
		t.Fatal(err)
	}
	if seats, err := service.AvailableSeats("matinee"); err != nil || seats != 1 {
		t.Fatalf("unexpected seats: %d %v", seats, err)
	}
	tickets := NewTicketService(db, nil)
	if _, err = tickets.PlaceOrder(model.TicketOrder{ID: "order-1", PerformanceID: "matinee", Buyer: "李", Quantity: 1}); err != nil {
		t.Fatal(err)
	}
	if total, err := tickets.QuantityByBuyer("李"); err != nil || total != 1 {
		t.Fatalf("buyer total failed: %d %v", total, err)
	}
	settlement, err := NewSettlementService(db, nil).Calculate("matinee")
	if err != nil || settlement.NetTickets != 4 {
		t.Fatalf("settlement failed: %+v %v", settlement, err)
	}
}
