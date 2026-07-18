package views

import (
	"time"

	"github.com/google/uuid"
)

// ViewData holds common data for all pages.
type ViewData struct {
	Title    string
	UserName string
	UserRole string
	TenantID uuid.UUID
	Path     string
}

// TicketData holds ticket-related view data.
type TicketData struct {
	ID           uuid.UUID
	Title        string
	Description  string
	Status       string
	Priority     string
	Category     string
	CreatedBy    uuid.UUID
	CreatorName  string
	AssignedTo   *uuid.UUID
	AssigneeName *string
	SLADeadline  *time.Time
	SLABreached  bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// CommentData holds comment view data.
type CommentData struct {
	ID         uuid.UUID
	TicketID   uuid.UUID
	UserID     uuid.UUID
	UserName   string
	Content    string
	IsInternal bool
	CreatedAt  time.Time
}

// HistoryData holds history entry view data.
type HistoryData struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	UserName  string
	Field     string
	OldValue  *string
	NewValue  *string
	CreatedAt time.Time
}

// ArticleData holds article view data.
type ArticleData struct {
	ID         uuid.UUID
	Title      string
	Content    string
	Category   string
	Tags       []string
	Published  bool
	CreatedBy  uuid.UUID
	AuthorName string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// UserData holds user view data.
type UserData struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Role      string
	AvatarURL *string
	Active    bool
	CreatedAt time.Time
}

// SLAPolicyData holds SLA policy view data.
type SLAPolicyData struct {
	ID                    uuid.UUID
	Priority              string
	ResponseTimeMinutes   int
	ResolutionTimeMinutes int
	EscalationMinutes     int
}

// AutomationData holds automation rule view data.
type AutomationData struct {
	ID          uuid.UUID
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
