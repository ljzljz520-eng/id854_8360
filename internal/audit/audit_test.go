package audit

import (
	"path/filepath"
	"testing"
	"theatrecontrol/internal/store"
)

func TestAuditSummary(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "audit.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	service := NewService(db)
	if err = service.Record("leader", "save", "role", "r1", "menus"); err != nil {
		t.Fatal(err)
	}
	if err = service.Record("stage", "open", "performance", "p1", "status"); err != nil {
		t.Fatal(err)
	}
	summary, err := service.Summary()
	if err != nil || summary != "2 entries [open:p1,save:r1]" {
		t.Fatalf("unexpected summary: %s %v", summary, err)
	}
}
