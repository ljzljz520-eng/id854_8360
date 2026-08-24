package rehearsal

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"theatrecontrol/internal/model"
)

type CrewRequirement struct {
	Role        string
	Minimum     int
	Description string
}

var defaultRequirements = []CrewRequirement{
	{Role: "actor", Minimum: 1, Description: "至少一名演员"},
	{Role: "stage_manager", Minimum: 1, Description: "至少一名舞监"},
}

func RequirementFor(role string) (CrewRequirement, bool) {
	for _, requirement := range defaultRequirements {
		if requirement.Role == role {
			return requirement, true
		}
	}
	return CrewRequirement{}, false
}

func ValidateCrew(participants []Participant) error {
	counts := map[string]int{}
	for _, participant := range participants {
		if strings.TrimSpace(participant.Role) == "" {
			return errors.New("participant role is empty")
		}
		counts[participant.Role]++
	}
	for _, requirement := range defaultRequirements {
		if counts[requirement.Role] < requirement.Minimum {
			return fmt.Errorf("crew needs %d %s", requirement.Minimum, requirement.Role)
		}
	}
	return nil
}

func SortParticipants(participants []Participant) []Participant {
	result := append([]Participant(nil), participants...)
	sort.Slice(result, func(i, j int) bool {
		if result[i].Role == result[j].Role {
			return result[i].ID < result[j].ID
		}
		return result[i].Role < result[j].Role
	})
	return result
}

func ParticipantRoles(participants []Participant) []string {
	seen := map[string]bool{}
	result := []string{}
	for _, participant := range participants {
		if !seen[participant.Role] {
			seen[participant.Role] = true
			result = append(result, participant.Role)
		}
	}
	sort.Strings(result)
	return result
}

func (s *RosterService) ValidateCrew(rehearsalID string) error {
	participants, err := s.List(rehearsalID)
	if err != nil {
		return err
	}
	return ValidateCrew(participants)
}

func (s *RosterService) RosterLabel(rehearsalID string) (string, error) {
	participants, err := s.List(rehearsalID)
	if err != nil {
		return "", err
	}
	participants = SortParticipants(participants)
	labels := make([]string, 0, len(participants))
	for _, participant := range participants {
		labels = append(labels, participant.Name+"/"+participant.Role)
	}
	return strings.Join(labels, ", "), nil
}

func ParticipantFromAssignment(value model.Assignment) Participant {
	return Participant{ID: value.ID, RehearsalID: value.ResourceID, PersonID: value.RoleID, Role: value.Kind, Attendance: value.Status}
}
