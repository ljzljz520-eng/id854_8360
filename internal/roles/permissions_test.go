package roles

import (
	"testing"
	"theatrecontrol/internal/model"
)

func TestSanitizeMenus(t *testing.T) {
	menus := sanitizeMenus([]string{"show", "bad", "show", "audit"})
	if len(menus) != 2 || menus[0] != "audit" || menus[1] != "show" {
		t.Fatalf("unexpected menus: %v", menus)
	}
}

func TestRolePermissionLifecycle(t *testing.T) {
	role := defaultActive(model.Role{ID: "r", Name: "舞监", UserType: model.UserStage, Menus: []string{"show"}})
	if !allows(role, "show") || allows(role, "ban") {
		t.Fatalf("permission evaluation failed")
	}
	role.Active = false
	if allows(role, "show") {
		t.Fatal("inactive role should not allow menu")
	}
}

func TestMenuPolicySuggestions(t *testing.T) {
	menus := model.SuggestedMenus(model.UserActor)
	if len(menus) != 2 || menus[0] != "costume" || menus[1] != "rehearsal" {
		t.Fatalf("unexpected suggestions: %v", menus)
	}
	if model.DisplayUserType(model.UserStage) != "舞监" || model.CanonicalMenuSet([]string{"show", "audit"}) != "audit,show" {
		t.Fatal("policy formatting failed")
	}
	if err := model.ValidateRoleShape(model.Role{ID: "r", Name: "演员", UserType: model.UserActor}); err != nil {
		t.Fatal(err)
	}
	if model.MenuSetDescription(nil) != "无可见权限" {
		t.Fatal("empty menu description failed")
	}
}
