package rehearsal

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"theatrecontrol/internal/audit"
	"theatrecontrol/internal/store"
)

type Participant struct {
	ID          string `json:"id"`
	RehearsalID string `json:"rehearsal_id"`
	PersonID    string `json:"person_id"`
	Name        string `json:"name"`
	Role        string `json:"role"`
	Attendance  string `json:"attendance"`
}

type RosterService struct {
	db    *store.DB
	audit *audit.Service
}

func NewRosterService(db *store.DB, logger *audit.Service) *RosterService {
	return &RosterService{db: db, audit: logger}
}

func validateParticipant(value Participant) error {
	if value.ID == "" || value.RehearsalID == "" || value.PersonID == "" {
		return errors.New("participant identity is incomplete")
	}
	if strings.TrimSpace(value.Name) == "" || strings.TrimSpace(value.Role) == "" {
		return errors.New("participant name and role are required")
	}
	return nil
}

func attendanceState(value string) string {
	if value == "" {
		return "expected"
	}
	if value == "expected" || value == "present" || value == "absent" || value == "excused" {
		return value
	}
	return "expected"
}

func (s *RosterService) Add(value Participant) (Participant, error) {
	if err := validateParticipant(value); err != nil {
		return Participant{}, err
	}
	if _, err := (&Service{db: s.db}).Get(value.RehearsalID); err != nil {
		return Participant{}, fmt.Errorf("rehearsal unavailable: %w", err)
	}
	value.Attendance = attendanceState(value.Attendance)
	if err := s.db.Put("assignment", value.ID, value); err != nil {
		return Participant{}, err
	}
	if s.audit != nil {
		_ = s.audit.Record(value.PersonID, "add_rehearsal_participant", "assignment", value.ID, value.RehearsalID)
	}
	return value, nil
}

func (s *RosterService) List(rehearsalID string) ([]Participant, error) {
	_, values, err := s.db.List("assignment")
	if err != nil {
		return nil, err
	}
	result := make([]Participant, 0, len(values))
	for _, data := range values {
		var value Participant
		if decodeParticipant(data, &value) == nil && value.RehearsalID == rehearsalID {
			result = append(result, value)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}

func (s *RosterService) MarkAttendance(id, attendance string) error {
	if attendanceState(attendance) != attendance {
		return errors.New("invalid attendance state")
	}
	var value Participant
	if err := s.db.Get("assignment", id, &value); err != nil {
		return err
	}
	value.Attendance = attendance
	if err := s.db.Put("assignment", id, value); err != nil {
		return err
	}
	if s.audit != nil {
		return s.audit.Record(value.PersonID, "mark_attendance", "assignment", id, attendance)
	}
	return nil
}

func (s *RosterService) AttendanceSummary(rehearsalID string) (map[string]int, error) {
	participants, err := s.List(rehearsalID)
	if err != nil {
		return nil, err
	}
	result := map[string]int{"expected": 0, "present": 0, "absent": 0, "excused": 0}
	for _, participant := range participants {
		result[participant.Attendance]++
	}
	return result, nil
}

func decodeParticipant(data []byte, target *Participant) error { return decodeRosterJSON(data, target) }
