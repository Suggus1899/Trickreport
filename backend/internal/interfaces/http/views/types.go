package views

import "time"

// ViewData holds common data for all pages.
type ViewData struct {
	Title    string
	UserName string
	UserRole string
	TenantID string
	Path     string
}

// TicketData holds ticket-related view data.
type TicketData struct {
	ID           string
	Title        string
	Description  string
	Status       string
	Priority     string
	Category     string
	CreatedBy    string
	CreatorName  string
	AssignedTo   *string
	AssigneeName *string
	SLADeadline  *time.Time
	SLABreached  bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// CommentData holds comment view data.
type CommentData struct {
	ID         string
	TicketID   string
	UserID     string
	UserName   string
	Content    string
	IsInternal bool
	CreatedAt  time.Time
}

// HistoryData holds history entry view data.
type HistoryData struct {
	ID        string
	UserID    string
	UserName  string
	Field     string
	OldValue  *string
	NewValue  *string
	CreatedAt time.Time
}

// ArticleData holds article view data.
type ArticleData struct {
	ID         string
	Title      string
	Content    string
	Category   string
	Tags       []string
	Published  bool
	CreatedBy  string
	AuthorName string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// UserData holds user view data.
type UserData struct {
	ID        string
	Name      string
	Email     string
	Role      string
	AvatarURL *string
	Active    bool
	CreatedAt time.Time
}

// SLAPolicyData holds SLA policy view data.
type SLAPolicyData struct {
	ID                    string
	Priority              string
	ResponseTimeMinutes   int
	ResolutionTimeMinutes int
	EscalationMinutes     int
}

// AutomationData holds automation rule view data.
type AutomationData struct {
	ID          string
	Name        string
	Description string
	TriggerType string
	IsActive    bool
	CreatedAt   time.Time
}

// SummaryData holds analytics summary view data.
type SummaryData struct {
	TotalTickets    int
	OpenTickets     int
	ResolvedTickets int
	SLABreached     int
}

// VolumePointData holds analytics volume point view data.
type VolumePointData struct {
	Date  string
	Count int
}

// StatusDistData holds analytics status distribution view data.
type StatusDistData struct {
	Status string
	Count  int
}

// ResolutionData holds analytics resolution time view data.
type ResolutionData struct {
	Priority string
	AvgHours float64
}
