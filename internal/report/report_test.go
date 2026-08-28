package report

import (
	"path/filepath"
	"strings"
	"testing"
	"theatrecontrol/internal/model"
	"theatrecontrol/internal/store"
)

func TestRenderText(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "report.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err = db.Put("role", "leader", model.Role{ID: "leader", Name: "团长", UserType: model.UserLeader, Active: true}); err != nil {
		t.Fatal(err)
	}
	text, err := NewService(db).RenderText()
	if err != nil {
		t.Fatal(err)
	}
	if len(text) == 0 {
		t.Fatal("empty report")
	}
	if matrix, err := NewService(db).RoleMatrix(); err != nil || !strings.Contains(matrix, "leader") {
		t.Fatalf("matrix failed: %s %v", matrix, err)
	}
	if snapshot, err := NewService(db).WorkflowSnapshot(); err != nil || snapshot["entities"] == "" {
		t.Fatalf("snapshot failed: %v", err)
	}
}
