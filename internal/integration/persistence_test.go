package integration

import (
	"path/filepath"
	"testing"

	"theatrecontrol/internal/model"
	"theatrecontrol/internal/roles"
	"theatrecontrol/internal/store"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "reopen.db")
	db, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	service := roles.NewService(db, nil)
	if _, err = service.SaveRole(model.Role{ID: "persist-role", Name: "团长", UserType: model.UserLeader, Menus: []string{"audit"}}); err != nil {
		t.Fatal(err)
	}
	if err = db.ApplyBatch([]store.Operation{{Kind: "assignment", ID: "batch-1", Value: model.Assignment{ID: "batch-1", RoleID: "persist-role", ResourceID: "room", Kind: "rehearsal", Status: model.StatusReady}}}); err != nil {
		t.Fatal(err)
	}
	if keys, err := db.Query(store.Query{Kind: "assignment", Prefix: "batch", Limit: 2}); err != nil || len(keys) != 1 {
		t.Fatalf("batch query failed: %v %v", keys, err)
	}
	if err = db.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	loaded, err := roles.NewService(reopened, nil).GetRole("persist-role")
	if err != nil || loaded.Name != "团长" || !model.ContainsMenu(loaded, "audit") {
		t.Fatalf("reopen lost role: %+v %v", loaded, err)
	}
}
