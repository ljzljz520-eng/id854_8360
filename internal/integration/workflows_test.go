package integration

import (
	"path/filepath"
	"testing"

	"theatrecontrol/internal/audit"
	"theatrecontrol/internal/model"
	"theatrecontrol/internal/rehearsal"
	"theatrecontrol/internal/report"
	"theatrecontrol/internal/roles"
	"theatrecontrol/internal/show"
	"theatrecontrol/internal/store"
	"theatrecontrol/internal/wardrobe"
)

func testDB(t *testing.T) *store.DB {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "control.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestWorkflowOne(t *testing.T) {
	db := testDB(t)
	logger := audit.NewService(db)
	roleService := roles.NewService(db, logger)
	role, err := roleService.SaveRole(model.Role{ID: "actor-role", Name: "演员", UserType: model.UserActor, Menus: []string{"rehearsal", "costume"}})
	if err != nil || !model.ContainsMenu(role, "rehearsal") {
		t.Fatalf("role setup failed: %v", err)
	}
	rehearsalService := rehearsal.NewService(db, logger)
	item, err := rehearsalService.Schedule(model.Rehearsal{ID: "r1", Production: "光", Room: "排练厅一", Leader: "stage", StartSlot: 10, EndSlot: 12, Status: model.StatusReady})
	if err != nil || item.Status != model.StatusReady {
		t.Fatalf("schedule failed: %v", err)
	}
	costumeService := wardrobe.NewService(db, logger)
	if _, err = costumeService.RegisterCostume(model.Costume{ID: "c1", Name: "红裙", Size: "M", Quantity: 2}); err != nil {
		t.Fatal(err)
	}
	if _, err = costumeService.Allocate(model.CostumeAllocation{ID: "a1", CostumeID: "c1", ActorID: "actor-1", Quantity: 1, Status: model.StatusReady}); err != nil {
		t.Fatal(err)
	}
	showService := show.NewService(db, logger)
	if _, err = showService.CreatePerformance(model.Performance{ID: "p1", Title: "光", Venue: "大剧场", Capacity: 100, Status: model.PerformancePlanned}); err != nil {
		t.Fatal(err)
	}
	if err = showService.Transition("p1", model.PerformanceOpen); err != nil {
		t.Fatal(err)
	}
	if _, err = showService.SellTickets("p1", 12); err != nil {
		t.Fatal(err)
	}
	summary, err := report.NewService(db).Summary()
	if err != nil || summary.ActiveRoles != 1 || summary.OpenPerformances != 1 {
		t.Fatalf("unexpected summary: %+v %v", summary, err)
	}
	banService := roles.NewBanService(db, logger)
	if _, err = banService.Save(model.BanRule{ID: "ban-1", RoleID: "actor-role", Menu: "ban", Reason: "演出期间关闭"}); err != nil {
		t.Fatal(err)
	}
	if banned, err := banService.IsBanned("actor-role", "ban"); err != nil || !banned {
		t.Fatalf("ban was not applied: %v", err)
	}
	ledger := wardrobe.NewLedgerService(db, logger)
	if _, err = ledger.Adjust(wardrobe.LedgerEntry{ID: "restock-1", CostumeID: "c1", Delta: 1, Reason: "返库"}); err != nil {
		t.Fatal(err)
	}
	if fit := wardrobe.CheckFit(model.Costume{ID: "c1", Size: "M"}, "L"); !fit.Approved {
		t.Fatalf("fit check failed: %+v", fit)
	}
}

func TestWorkflowTwo(t *testing.T) {
	db := testDB(t)
	logger := audit.NewService(db)
	roleService := roles.NewService(db, logger)
	if _, err := roleService.SaveRole(model.Role{ID: "ticket-role", Name: "票务", UserType: model.UserTicket, Menus: []string{"show"}}); err != nil {
		t.Fatal(err)
	}
	showService := show.NewService(db, logger)
	if _, err := showService.CreatePerformance(model.Performance{ID: "p2", Title: "夜航", Venue: "黑匣子", Capacity: 30, Status: model.PerformancePlanned}); err != nil {
		t.Fatal(err)
	}
	if err := showService.Transition("p2", model.PerformanceOpen); err != nil {
		t.Fatal(err)
	}
	if _, err := showService.SellTickets("p2", 30); err != nil {
		t.Fatal(err)
	}
	if _, err := showService.SellTickets("p2", 1); err == nil {
		t.Fatal("expected capacity error")
	}
	if err := showService.Transition("p2", model.PerformanceClosed); err != nil {
		t.Fatal(err)
	}
	if available, err := showService.AvailableSeats("p2"); err != nil || available != 0 {
		t.Fatalf("unexpected seats: %d %v", available, err)
	}
	tickets := show.NewTicketService(db, logger)
	if _, err := tickets.PlaceOrder(model.TicketOrder{ID: "order-2", PerformanceID: "p2", Buyer: "票务组", Quantity: 1}); err == nil {
		t.Fatal("closed show accepted order")
	}
	settlement, err := show.NewSettlementService(db, logger).Calculate("p2")
	if err != nil || settlement.NetTickets != 30 {
		t.Fatalf("unexpected settlement: %+v %v", settlement, err)
	}
}

func TestWorkflowThree(t *testing.T) {
	db := testDB(t)
	logger := audit.NewService(db)
	roleService := roles.NewService(db, logger)
	if _, err := roleService.SaveRole(model.Role{ID: "stage-role", Name: "舞监", UserType: model.UserStage, Menus: []string{"rehearsal", "show", "ban"}}); err != nil {
		t.Fatal(err)
	}
	if ok, err := roleService.Allows("stage-role", "ban"); err != nil || !ok {
		t.Fatalf("permission missing: %v", err)
	}
	rehearsalService := rehearsal.NewService(db, logger)
	if _, err := rehearsalService.Schedule(model.Rehearsal{ID: "r3", Production: "远方", Room: "排练厅二", Leader: "stage", StartSlot: 8, EndSlot: 9}); err != nil {
		t.Fatal(err)
	}
	roster := rehearsal.NewRosterService(db, logger)
	if _, err := roster.Add(rehearsal.Participant{ID: "crew-1", RehearsalID: "r3", PersonID: "actor-1", Name: "林", Role: "actor"}); err != nil {
		t.Fatal(err)
	}
	if _, err := roster.Add(rehearsal.Participant{ID: "crew-2", RehearsalID: "r3", PersonID: "stage-1", Name: "周", Role: "stage_manager"}); err != nil {
		t.Fatal(err)
	}
	if err := roster.ValidateCrew("r3"); err != nil {
		t.Fatal(err)
	}
	if err := roster.MarkAttendance("crew-1", "present"); err != nil {
		t.Fatal(err)
	}
	if _, err := rehearsalService.Schedule(model.Rehearsal{ID: "r4", Production: "远方", Room: "排练厅二", Leader: "stage", StartSlot: 8, EndSlot: 10}); err == nil {
		t.Fatal("expected overlap error")
	}
	if err := rehearsalService.Transition("r3", model.StatusComplete); err != nil {
		t.Fatal(err)
	}
	text, err := report.NewService(db).RenderText()
	if err != nil || len(text) == 0 {
		t.Fatalf("report failed: %v", err)
	}
	if label, err := roster.RosterLabel("r3"); err != nil || label == "" {
		t.Fatalf("roster label failed: %s %v", label, err)
	}
}
