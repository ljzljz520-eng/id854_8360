package rehearsal

import (
	"path/filepath"
	"testing"
	"theatrecontrol/internal/model"
	"theatrecontrol/internal/store"
)

func TestScheduleRejectsOverlap(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "rehearsal.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	service := NewService(db, nil)
	if _, err = service.Schedule(model.Rehearsal{ID: "one", Production: "戏", Room: "一号厅", StartSlot: 9, EndSlot: 11}); err != nil {
		t.Fatal(err)
	}
	if _, err = service.Schedule(model.Rehearsal{ID: "two", Production: "戏", Room: "一号厅", StartSlot: 10, EndSlot: 12}); err == nil {
		t.Fatal("expected overlap")
	}
}

func TestScheduleTransition(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "transition.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	service := NewService(db, nil)
	if _, err = service.Schedule(model.Rehearsal{ID: "three", Production: "戏", Room: "二号厅", StartSlot: 13, EndSlot: 15}); err != nil {
		t.Fatal(err)
	}
	if err = service.Transition("three", model.StatusActive); err != nil {
		t.Fatal(err)
	}
	item, err := service.Get("three")
	if err != nil || item.Status != model.StatusActive {
		t.Fatalf("bad transition: %+v %v", item, err)
	}
}
