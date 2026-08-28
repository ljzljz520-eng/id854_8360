package roles

import (
	"path/filepath"
	"testing"

	"theatrecontrol/internal/model"
	"theatrecontrol/internal/store"
)

func TestRoleAllowsEmptyMenuSelection(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "role.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	service := NewService(db, nil)
	role, err := service.SaveRole(model.Role{ID: "empty-role", Name: "演员", UserType: model.UserActor, Menus: []string{}})
	if err != nil {
		t.Fatal(err)
	}
	if role.ID != "empty-role" || len(role.Menus) != 0 {
		t.Fatalf("unexpected role: %+v", role)
	}
	loaded, err := service.GetRole("empty-role")
	if err != nil || len(loaded.Menus) != 0 {
		t.Fatalf("empty role was not persisted: %+v %v", loaded, err)
	}
}
