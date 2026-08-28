package model

type UserType string

const (
	UserLeader UserType = "leader"
	UserActor  UserType = "actor"
	UserStage  UserType = "stage_manager"
	UserTicket UserType = "ticketing"
)

type Role struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	UserType    UserType `json:"user_type"`
	Menus       []string `json:"menus"`
	Active      bool     `json:"active"`
	Description string   `json:"description"`
}

type Rehearsal struct {
	ID         string `json:"id"`
	Production string `json:"production"`
	Room       string `json:"room"`
	Leader     string `json:"leader"`
	Status     string `json:"status"`
	StartSlot  int    `json:"start_slot"`
	EndSlot    int    `json:"end_slot"`
}

type Costume struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Size     string `json:"size"`
	Quantity int    `json:"quantity"`
	Status   string `json:"status"`
}

type CostumeAllocation struct {
	ID        string `json:"id"`
	CostumeID string `json:"costume_id"`
	ActorID   string `json:"actor_id"`
	Quantity  int    `json:"quantity"`
	Status    string `json:"status"`
}

type Performance struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Venue    string `json:"venue"`
	Status   string `json:"status"`
	Capacity int    `json:"capacity"`
	Sold     int    `json:"sold"`
}

type BanRule struct {
	ID        string   `json:"id"`
	RoleID    string   `json:"role_id"`
	Menu      string   `json:"menu"`
	Reason    string   `json:"reason"`
	Enabled   bool     `json:"enabled"`
	AuditTags []string `json:"audit_tags"`
}

type AuditEntry struct {
	ID        string `json:"id"`
	Actor     string `json:"actor"`
	Action    string `json:"action"`
	Entity    string `json:"entity"`
	EntityID  string `json:"entity_id"`
	Timestamp int64  `json:"timestamp"`
	Details   string `json:"details"`
}

type Assignment struct {
	ID         string `json:"id"`
	RoleID     string `json:"role_id"`
	ResourceID string `json:"resource_id"`
	Kind       string `json:"kind"`
	Status     string `json:"status"`
}

type TicketOrder struct {
	ID            string `json:"id"`
	PerformanceID string `json:"performance_id"`
	Buyer         string `json:"buyer"`
	Quantity      int    `json:"quantity"`
	Status        string `json:"status"`
	Note          string `json:"note"`
}
