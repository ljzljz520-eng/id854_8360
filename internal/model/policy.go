package model

import (
	"errors"
	"sort"
	"strings"
)

type MenuPolicy struct {
	Code        string
	Label       string
	Audience    []UserType
	Description string
}

var MenuCatalog = []MenuPolicy{
	{Code: "rehearsal", Label: "排练", Audience: []UserType{UserLeader, UserActor, UserStage}},
	{Code: "costume", Label: "服装", Audience: []UserType{UserLeader, UserActor, UserStage}},
	{Code: "show", Label: "演出", Audience: []UserType{UserLeader, UserStage, UserTicket}},
	{Code: "ban", Label: "封禁", Audience: []UserType{UserLeader, UserStage}},
	{Code: "audit", Label: "日志", Audience: []UserType{UserLeader, UserStage, UserTicket}},
}

func MenuPolicyFor(code string) (MenuPolicy, bool) {
	for _, policy := range MenuCatalog {
		if policy.Code == code {
			return policy, true
		}
	}
	return MenuPolicy{}, false
}

func AudienceIncludes(policy MenuPolicy, userType UserType) bool {
	for _, audience := range policy.Audience {
		if audience == userType {
			return true
		}
	}
	return false
}

func SuggestedMenus(userType UserType) []string {
	result := make([]string, 0)
	for _, policy := range MenuCatalog {
		if AudienceIncludes(policy, userType) {
			result = append(result, policy.Code)
		}
	}
	sort.Strings(result)
	return result
}

func ValidateEntityID(value string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return errors.New("entity id is empty")
	}
	if len(trimmed) > 64 {
		return errors.New("entity id is too long")
	}
	for _, r := range trimmed {
		if r == '/' || r == '\\' || r == ' ' {
			return errors.New("entity id contains a forbidden separator")
		}
	}
	return nil
}

func CompareMenus(left, right []string) (same, onlyLeft, onlyRight []string) {
	leftSet, rightSet := map[string]bool{}, map[string]bool{}
	for _, item := range left {
		leftSet[item] = true
	}
	for _, item := range right {
		rightSet[item] = true
	}
	for item := range leftSet {
		if rightSet[item] {
			same = append(same, item)
		} else {
			onlyLeft = append(onlyLeft, item)
		}
	}
	for item := range rightSet {
		if !leftSet[item] {
			onlyRight = append(onlyRight, item)
		}
	}
	sort.Strings(same)
	sort.Strings(onlyLeft)
	sort.Strings(onlyRight)
	return
}

func MergeMenuSelections(primary, fallback []string) []string {
	seen := map[string]bool{}
	result := make([]string, 0, len(primary)+len(fallback))
	for _, list := range [][]string{primary, fallback} {
		for _, menu := range list {
			if menu != "" && !seen[menu] {
				seen[menu] = true
				result = append(result, menu)
			}
		}
	}
	sort.Strings(result)
	return result
}
