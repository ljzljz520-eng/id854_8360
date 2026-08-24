package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"theatrecontrol/internal/audit"
	"theatrecontrol/internal/model"
	"theatrecontrol/internal/report"
	"theatrecontrol/internal/roles"
	"theatrecontrol/internal/show"
	"theatrecontrol/internal/store"
)

func main() {
	path := filepath.Join(os.TempDir(), "theatre-controlroom.db")
	db, err := store.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	roleService := roles.NewService(db, audit.NewService(db))
	role, err := roleService.SaveRole(model.Role{ID: "demo-actor", Name: "演员", UserType: model.UserActor, Menus: []string{"rehearsal", "show"}})
	if err != nil {
		fmt.Println(err)
		return
	}
	showService := show.NewService(db, audit.NewService(db))
	_, _ = showService.CreatePerformance(model.Performance{ID: "demo-show", Title: "城市之光", Venue: "主舞台", Status: model.PerformancePlanned})
	reporter := report.NewService(db)
	summary, _ := reporter.Summary()
	output := map[string]any{"role": role, "summary": summary}
	encoded, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println(string(encoded))
}
