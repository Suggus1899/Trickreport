package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	appCX "github.com/trickreport/backend/internal/application/cx"
	domain "github.com/trickreport/backend/internal/domain/cx"
	"github.com/trickreport/backend/internal/interfaces/http/middleware"
	"github.com/trickreport/backend/internal/interfaces/http/response"
	"github.com/trickreport/backend/internal/interfaces/http/validator"
)

// CXHandler serves the Stage 5 modules: external customers and
// contacts, the inbox across portal, email, live chat and WhatsApp,
// assignment, SLA, escalations, channel configuration, assignment
// rules, auto-replies, surveys and reports. Write endpoints for
// staff-only resources are additionally guarded by RequireRole in
// server.go; the per-customer isolation rule itself is enforced in
// the application service, so an external customer can never read
// another customer's conversations even if the UI showed them.
// Internal notes are filtered in the service for external customers.
type CXHandler struct {
	svc *appCX.Service
}

// NewCXHandler creates a CXHandler.
func NewCXHandler(svc *appCX.Service) *CXHandler { return &CXHandler{svc: svc} }

func cxErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		response.Error(w, http.StatusNotFound, "not found")
	case errors.Is(err, domain.ErrValidation):
		response.Error(w, http.StatusBadRequest, "validation failed")
	case errors.Is(err, domain.ErrForbidden):
		response.Error(w, http.StatusForbidden, "insufficient permissions")
	case errors.Is(err, domain.ErrInvalidTransition):
		response.Error(w, http.StatusConflict, "invalid status transition")
	case errors.Is(err, domain.ErrConflict):
		response.Error(w, http.StatusConflict, "conflict")
	default:
		response.Error(w, http.StatusInternalServerError, "internal error")
	}
}

func cxTenant(r *http.Request) uuid.UUID { return middleware.TenantFromContext(r.Context()) }

func cxCaller(r *http.Request) (uuid.UUID, string) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		return uuid.Nil, ""
	}
	return claims.UserID, claims.Role
}

func cxID(r *http.Request, name string) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, name))
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}

func cxOptID(s string) *uuid.UUID {
	if s == "" {
		return nil
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return nil
	}
	return &id
}

func cxDecode(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return false
	}
	if err := validator.Validate.Struct(dst); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{"error": "validation failed", "details": validator.ValidationError(err)})
		return false
	}
	return true
}

// ── DTOs ────────────────────────────────────────────────────────────

type cxCustomerDTO struct {
	ID               uuid.UUID `json:"id"`
	TenantID         uuid.UUID `json:"tenant_id"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	Phone            string    `json:"phone"`
	Company          string    `json:"company"`
	Status           string    `json:"status"`
	ConsentGiven     bool      `json:"consent_given"`
	ErasureRequested bool      `json:"erasure_requested"`
	CreatedAt        string    `json:"created_at"`
	UpdatedAt        string    `json:"updated_at"`
}

func toCXCustomerDTO(c *domain.Customer) cxCustomerDTO {
	return cxCustomerDTO{ID: c.ID, TenantID: c.TenantID, Name: c.Name, Email: c.Email, Phone: c.Phone, Company: c.Company, Status: c.Status, ConsentGiven: c.ConsentGiven, ErasureRequested: c.ErasureRequested, CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"), UpdatedAt: c.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")}
}

type cxCustomerInput struct {
	Name             string `json:"name" validate:"required,max=255"`
	Email            string `json:"email" validate:"max=255"`
	Phone            string `json:"phone" validate:"max=100"`
	Company          string `json:"company" validate:"max=255"`
	Status           string `json:"status"`
	ConsentGiven     bool   `json:"consent_given"`
	ErasureRequested bool   `json:"erasure_requested"`
}

type cxContactDTO struct {
	ID         uuid.UUID  `json:"id"`
	TenantID   uuid.UUID  `json:"tenant_id"`
	CustomerID uuid.UUID  `json:"customer_id"`
	UserID     *uuid.UUID `json:"user_id"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	Role       string     `json:"role"`
	CreatedAt  string     `json:"created_at"`
}

func toCXContactDTO(c *domain.CustomerContact) cxContactDTO {
	return cxContactDTO{ID: c.ID, TenantID: c.TenantID, CustomerID: c.CustomerID, UserID: c.UserID, Name: c.Name, Email: c.Email, Role: c.Role, CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z07:00")}
}

type cxContactInput struct {
	UserID *uuid.UUID `json:"user_id"`
	Name   string     `json:"name" validate:"required,max=255"`
	Email  string     `json:"email" validate:"required,max=255"`
	Role   string     `json:"role"`
}

type cxConversationDTO struct {
	ID                   uuid.UUID  `json:"id"`
	TenantID             uuid.UUID  `json:"tenant_id"`
	CustomerID           uuid.UUID  `json:"customer_id"`
	CustomerName         string     `json:"customer_name"`
	Channel              string     `json:"channel"`
	Subject              string     `json:"subject"`
	Status               string     `json:"status"`
	Priority             string     `json:"priority"`
	AssignedTo           *uuid.UUID `json:"assigned_to"`
	AssigneeName         string     `json:"assignee_name"`
	CreatedBy            *uuid.UUID `json:"created_by"`
	ExternalThreadID     string     `json:"external_thread_id"`
	SLAResponseMinutes   int        `json:"sla_response_minutes"`
	SLAResolutionMinutes int        `json:"sla_resolution_minutes"`
	ResponseDueAt        *string    `json:"response_due_at"`
	ResolutionDueAt      *string    `json:"resolution_due_at"`
	FirstResponseAt      *string    `json:"first_response_at"`
	Escalated            bool       `json:"escalated"`
	CreatedAt            string     `json:"created_at"`
	UpdatedAt            string     `json:"updated_at"`
}

func cxFormatTime(t interface{ Format(string) string }) *string {
	s := t.Format("2006-01-02T15:04:05Z07:00")
	return &s
}

func toCXConversationDTO(c *domain.Conversation) cxConversationDTO {
	dto := cxConversationDTO{ID: c.ID, TenantID: c.TenantID, CustomerID: c.CustomerID, CustomerName: c.CustomerName, Channel: c.Channel, Subject: c.Subject, Status: c.Status, Priority: c.Priority, AssignedTo: c.AssignedTo, AssigneeName: c.AssigneeName, CreatedBy: c.CreatedBy, ExternalThreadID: c.ExternalThreadID, SLAResponseMinutes: c.SLAResponseMinutes, SLAResolutionMinutes: c.SLAResolutionMinutes, Escalated: c.Escalated, CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"), UpdatedAt: c.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")}
	if c.ResponseDueAt != nil {
		dto.ResponseDueAt = cxFormatTime(c.ResponseDueAt)
	}
	if c.ResolutionDueAt != nil {
		dto.ResolutionDueAt = cxFormatTime(c.ResolutionDueAt)
	}
	if c.FirstResponseAt != nil {
		dto.FirstResponseAt = cxFormatTime(c.FirstResponseAt)
	}
	return dto
}

type cxConversationInput struct {
	CustomerID       uuid.UUID `json:"customer_id" validate:"required"`
	Channel          string    `json:"channel" validate:"required"`
	Subject          string    `json:"subject" validate:"required,max=255"`
	Priority         string    `json:"priority"`
	ExternalThreadID string    `json:"external_thread_id" validate:"max=255"`
}

type cxStatusInput struct {
	Status string `json:"status" validate:"required"`
}

type cxMessageDTO struct {
	ID                uuid.UUID  `json:"id"`
	TenantID          uuid.UUID  `json:"tenant_id"`
	ConversationID    uuid.UUID  `json:"conversation_id"`
	SenderType        string     `json:"sender_type"`
	SenderID          *uuid.UUID `json:"sender_id"`
	SenderName        string     `json:"sender_name"`
	Body              string     `json:"body"`
	IsInternal        bool       `json:"is_internal"`
	Channel           string     `json:"channel"`
	ExternalMessageID string     `json:"external_message_id"`
	CreatedAt         string     `json:"created_at"`
}

func toCXMessageDTO(m *domain.Message) cxMessageDTO {
	return cxMessageDTO{ID: m.ID, TenantID: m.TenantID, ConversationID: m.ConversationID, SenderType: m.SenderType, SenderID: m.SenderID, SenderName: m.SenderName, Body: m.Body, IsInternal: m.IsInternal, Channel: m.Channel, ExternalMessageID: m.ExternalMessageID, CreatedAt: m.CreatedAt.Format("2006-01-02T15:04:05Z07:00")}
}

type cxMessageInput struct {
	SenderType string `json:"sender_type"`
	Body       string `json:"body" validate:"required,max=20000"`
	IsInternal bool   `json:"is_internal"`
}

type cxEmailAccountDTO struct {
	ID           uuid.UUID `json:"id"`
	TenantID     uuid.UUID `json:"tenant_id"`
	Name         string    `json:"name"`
	SMTPHost     string    `json:"smtp_host"`
	SMTPPort     int       `json:"smtp_port"`
	SMTPUsername string    `json:"smtp_username"`
	SMTPFrom     string    `json:"smtp_from"`
	IMAPHost     string    `json:"imap_host"`
	IMAPPort     int       `json:"imap_port"`
	IMAPUsername string    `json:"imap_username"`
	Active       bool      `json:"active"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

func toCXEmailAccountDTO(a *domain.EmailAccount) cxEmailAccountDTO {
	return cxEmailAccountDTO{ID: a.ID, TenantID: a.TenantID, Name: a.Name, SMTPHost: a.SMTPHost, SMTPPort: a.SMTPPort, SMTPUsername: a.SMTPUsername, SMTPFrom: a.SMTPFrom, IMAPHost: a.IMAPHost, IMAPPort: a.IMAPPort, IMAPUsername: a.IMAPUsername, Active: a.Active, CreatedAt: a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"), UpdatedAt: a.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")}
}

type cxEmailAccountInput struct {
	Name         string `json:"name" validate:"required,max=255"`
	SMTPHost     string `json:"smtp_host" validate:"max=500"`
	SMTPPort     int    `json:"smtp_port"`
	SMTPUsername string `json:"smtp_username" validate:"max=255"`
	SMTPFrom     string `json:"smtp_from" validate:"max=255"`
	IMAPHost     string `json:"imap_host" validate:"max=500"`
	IMAPPort     int    `json:"imap_port"`
	IMAPUsername string `json:"imap_username" validate:"max=255"`
	Active       bool   `json:"active"`
}

type cxWhatsAppAccountDTO struct {
	ID            uuid.UUID `json:"id"`
	TenantID      uuid.UUID `json:"tenant_id"`
	Name          string    `json:"name"`
	BusinessID    string    `json:"business_id"`
	PhoneNumberID string    `json:"phone_number_id"`
	DisplayName   string    `json:"display_name"`
	Active        bool      `json:"active"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
}

func toCXWhatsAppAccountDTO(a *domain.WhatsAppAccount) cxWhatsAppAccountDTO {
	return cxWhatsAppAccountDTO{ID: a.ID, TenantID: a.TenantID, Name: a.Name, BusinessID: a.BusinessID, PhoneNumberID: a.PhoneNumberID, DisplayName: a.DisplayName, Active: a.Active, CreatedAt: a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"), UpdatedAt: a.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")}
}

type cxWhatsAppAccountInput struct {
	Name          string `json:"name" validate:"required,max=255"`
	BusinessID    string `json:"business_id" validate:"max=255"`
	PhoneNumberID string `json:"phone_number_id" validate:"max=255"`
	DisplayName   string `json:"display_name" validate:"max=255"`
	VerifyToken   string `json:"verify_token" validate:"max=255"`
	Active        bool   `json:"active"`
}

type cxSLAPolicyDTO struct {
	ID                uuid.UUID `json:"id"`
	TenantID          uuid.UUID `json:"tenant_id"`
	Channel           string    `json:"channel"`
	Priority          string    `json:"priority"`
	ResponseMinutes   int       `json:"response_minutes"`
	ResolutionMinutes int       `json:"resolution_minutes"`
	CreatedAt         string    `json:"created_at"`
	UpdatedAt         string    `json:"updated_at"`
}

func toCXSLAPolicyDTO(p *domain.SLAPolicy) cxSLAPolicyDTO {
	return cxSLAPolicyDTO{ID: p.ID, TenantID: p.TenantID, Channel: p.Channel, Priority: p.Priority, ResponseMinutes: p.ResponseMinutes, ResolutionMinutes: p.ResolutionMinutes, CreatedAt: p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"), UpdatedAt: p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")}
}

type cxSLAPolicyInput struct {
	Channel           string `json:"channel" validate:"required"`
	Priority          string `json:"priority" validate:"required"`
	ResponseMinutes   int    `json:"response_minutes"`
	ResolutionMinutes int    `json:"resolution_minutes"`
}

type cxAssignmentRuleDTO struct {
	ID             uuid.UUID `json:"id"`
	TenantID       uuid.UUID `json:"tenant_id"`
	Channel        string    `json:"channel"`
	Priority       string    `json:"priority"`
	AssigneeUserID uuid.UUID `json:"assignee_user_id"`
	Position       int       `json:"position"`
	Active         bool      `json:"active"`
	CreatedAt      string    `json:"created_at"`
}

func toCXAssignmentRuleDTO(r *domain.AssignmentRule) cxAssignmentRuleDTO {
	return cxAssignmentRuleDTO{ID: r.ID, TenantID: r.TenantID, Channel: r.Channel, Priority: r.Priority, AssigneeUserID: r.AssigneeUserID, Position: r.Position, Active: r.Active, CreatedAt: r.CreatedAt.Format("2006-01-02T15:04:05Z07:00")}
}

type cxAssignmentRuleInput struct {
	Channel        string    `json:"channel"`
	Priority       string    `json:"priority"`
	AssigneeUserID uuid.UUID `json:"assignee_user_id" validate:"required"`
	Position       int       `json:"position"`
	Active         bool      `json:"active"`
}

type cxAutoReplyDTO struct {
	ID           uuid.UUID `json:"id"`
	TenantID     uuid.UUID `json:"tenant_id"`
	TriggerKind  string    `json:"trigger_kind"`
	Keyword      string    `json:"keyword"`
	Channel      string    `json:"channel"`
	ResponseText string    `json:"response_text"`
	Active       bool      `json:"active"`
	CreatedAt    string    `json:"created_at"`
}

func toCXAutoReplyDTO(a *domain.AutoReply) cxAutoReplyDTO {
	return cxAutoReplyDTO{ID: a.ID, TenantID: a.TenantID, TriggerKind: a.TriggerKind, Keyword: a.Keyword, Channel: a.Channel, ResponseText: a.ResponseText, Active: a.Active, CreatedAt: a.CreatedAt.Format("2006-01-02T15:04:05Z07:00")}
}

type cxAutoReplyInput struct {
	TriggerKind  string `json:"trigger_kind" validate:"required"`
	Keyword      string `json:"keyword" validate:"max=255"`
	Channel      string `json:"channel"`
	ResponseText string `json:"response_text" validate:"required,max=5000"`
	Active       bool   `json:"active"`
}

type cxSurveyDTO struct {
	ID             uuid.UUID `json:"id"`
	TenantID       uuid.UUID `json:"tenant_id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	CustomerID     uuid.UUID `json:"customer_id"`
	Rating         int       `json:"rating"`
	Comment        string    `json:"comment"`
	CreatedAt      string    `json:"created_at"`
}

func toCXSurveyDTO(s *domain.Survey) cxSurveyDTO {
	return cxSurveyDTO{ID: s.ID, TenantID: s.TenantID, ConversationID: s.ConversationID, CustomerID: s.CustomerID, Rating: s.Rating, Comment: s.Comment, CreatedAt: s.CreatedAt.Format("2006-01-02T15:04:05Z07:00")}
}

type cxSurveyInput struct {
	Rating  int    `json:"rating" validate:"required,min=1,max=5"`
	Comment string `json:"comment" validate:"max=5000"`
}

type cxEmailIngestInput struct {
	AccountID         uuid.UUID `json:"account_id"`
	FromEmail         string    `json:"from_email" validate:"required,max=255"`
	Subject           string    `json:"subject" validate:"max=255"`
	Body              string    `json:"body" validate:"required,max=20000"`
	ThreadID          string    `json:"thread_id" validate:"max=255"`
	ExternalMessageID string    `json:"external_message_id" validate:"max=255"`
}

// ── Customers ───────────────────────────────────────────────────────

func (h *CXHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	viewerID, role := cxCaller(r)
	customers, err := h.svc.ListCustomers(r.Context(), cxTenant(r), viewerID, role)
	if err != nil {
		cxErr(w, err)
		return
	}
	out := make([]cxCustomerDTO, 0, len(customers))
	for i := range customers {
		out = append(out, toCXCustomerDTO(&customers[i]))
	}
	response.JSON(w, http.StatusOK, out)
}

func (h *CXHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var in cxCustomerInput
	if !cxDecode(w, r, &in) {
		return
	}
	c, err := h.svc.CreateCustomer(r.Context(), cxTenant(r), &domain.Customer{Name: in.Name, Email: in.Email, Phone: in.Phone, Company: in.Company, Status: in.Status, ConsentGiven: in.ConsentGiven, ErasureRequested: in.ErasureRequested})
	if err != nil {
		cxErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toCXCustomerDTO(c))
}

func (h *CXHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	id, ok := cxID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	viewerID, role := cxCaller(r)
	c, err := h.svc.GetCustomer(r.Context(), cxTenant(r), viewerID, role, id)
	if err != nil {
		cxErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toCXCustomerDTO(c))
}

func (h *CXHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	id, ok := cxID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in cxCustomerInput
	if !cxDecode(w, r, &in) {
		return
	}
	c, err := h.svc.UpdateCustomer(r.Context(), cxTenant(r), &domain.Customer{ID: id, Name: in.Name, Email: in.Email, Phone: in.Phone, Company: in.Company, Status: in.Status, ConsentGiven: in.ConsentGiven, ErasureRequested: in.ErasureRequested})
	if err != nil {
		cxErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toCXCustomerDTO(c))
}

func (h *CXHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	id, ok := cxID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteCustomer(r.Context(), cxTenant(r), id); err != nil {
		cxErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CXHandler) ListContacts(w http.ResponseWriter, r *http.Request) {
	id, ok := cxID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	viewerID, role := cxCaller(r)
	contacts, err := h.svc.ListContacts(r.Context(), cxTenant(r), viewerID, role, id)
	if err != nil {
		cxErr(w, err)
		return
	}
	out := make([]cxContactDTO, 0, len(contacts))
	for i := range contacts {
		out = append(out, toCXContactDTO(&contacts[i]))
	}
	response.JSON(w, http.StatusOK, out)
}

func (h *CXHandler) CreateContact(w http.ResponseWriter, r *http.Request) {
	id, ok := cxID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in cxContactInput
	if !cxDecode(w, r, &in) {
		return
	}
	c, err := h.svc.CreateContact(r.Context(), cxTenant(r), id, &domain.CustomerContact{UserID: in.UserID, Name: in.Name, Email: in.Email, Role: in.Role})
	if err != nil {
		cxErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toCXContactDTO(c))
}

func (h *CXHandler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	id, ok := cxID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	contactID, ok := cxID(r, "contactId")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteContact(r.Context(), cxTenant(r), id, contactID); err != nil {
		cxErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ── Conversations and messages ──────────────────────────────────────

func (h *CXHandler) ListConversations(w http.ResponseWriter, r *http.Request) {
	viewerID, role := cxCaller(r)
	convs, err := h.svc.ListConversations(r.Context(), cxTenant(r), viewerID, role, cxOptID(r.URL.Query().Get("customer_id")), r.URL.Query().Get("channel"), r.URL.Query().Get("status"), cxOptID(r.URL.Query().Get("assigned_to")))
	if err != nil {
		cxErr(w, err)
		return
	}
	out := make([]cxConversationDTO, 0, len(convs))
	for i := range convs {
		out = append(out, toCXConversationDTO(&convs[i]))
	}
	response.JSON(w, http.StatusOK, out)
}

func (h *CXHandler) CreateConversation(w http.ResponseWriter, r *http.Request) {
	var in cxConversationInput
	if !cxDecode(w, r, &in) {
		return
	}
	viewerID, role := cxCaller(r)
	c, err := h.svc.CreateConversation(r.Context(), cxTenant(r), viewerID, role, &domain.Conversation{CustomerID: in.CustomerID, Channel: in.Channel, Subject: in.Subject, Priority: in.Priority, ExternalThreadID: in.ExternalThreadID})
	if err != nil {
		cxErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toCXConversationDTO(c))
}

func (h *CXHandler) GetConversation(w http.ResponseWriter, r *http.Request) {
	id, ok := cxID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	viewerID, role := cxCaller(r)
	c, err := h.svc.GetConversation(r.Context(), cxTenant(r), viewerID, role, id)
	if err != nil {
		cxErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toCXConversationDTO(c))
}

func (h *CXHandler) UpdateConversationStatus(w http.ResponseWriter, r *http.Request) {
	id, ok := cxID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in cxStatusInput
	if !cxDecode(w, r, &in) {
		return
	}
	c, err := h.svc.UpdateConversationStatus(r.Context(), cxTenant(r), &domain.Conversation{ID: id}, in.Status)
	if err != nil {
		cxErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toCXConversationDTO(c))
}

func (h *CXHandler) TakeConversation(w http.ResponseWriter, r *http.Request) {
	id, ok := cxID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	callerID, _ := cxCaller(r)
	c, err := h.svc.TakeConversation(r.Context(), cxTenant(r), callerID, id)
	if err != nil {
		cxErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toCXConversationDTO(c))
}

func (h *CXHandler) ReleaseConversation(w http.ResponseWriter, r *http.Request) {
	id, ok := cxID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	callerID, role := cxCaller(r)
	c, err := h.svc.ReleaseConversation(r.Context(), cxTenant(r), callerID, role, id)
	if err != nil {
		cxErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toCXConversationDTO(c))
}

func (h *CXHandler) EscalateConversation(w http.ResponseWriter, r *http.Request) {
	id, ok := cxID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	c, err := h.svc.EscalateConversation(r.Context(), cxTenant(r), id)
	if err != nil {
		cxErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toCXConversationDTO(c))
}

func (h *CXHandler) ListMessages(w http.ResponseWriter, r *http.Request) {
	id, ok := cxID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	viewerID, role := cxCaller(r)
	messages, err := h.svc.ListMessages(r.Context(), cxTenant(r), viewerID, role, id)
	if err != nil {
		cxErr(w, err)
		return
	}
	out := make([]cxMessageDTO, 0, len(messages))
	for i := range messages {
		out = append(out, toCXMessageDTO(&messages[i]))
	}
	response.JSON(w, http.StatusOK, out)
}

func (h *CXHandler) AddMessage(w http.ResponseWriter, r *http.Request) {
	id, ok := cxID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in cxMessageInput
	if !cxDecode(w, r, &in) {
		return
	}
	viewerID, role := cxCaller(r)
	m, err := h.svc.AddMessage(r.Context(), cxTenant(r), viewerID, role, id, &domain.Message{SenderType: in.SenderType, Body: in.Body, IsInternal: in.IsInternal})
	if err != nil {
		cxErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toCXMessageDTO(m))
}

// ── Channel configuration and ingestion ─────────────────────────────

func (h *CXHandler) ListEmailAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.svc.ListEmailAccounts(r.Context(), cxTenant(r))
	if err != nil {
		cxErr(w, err)
		return
	}
	out := make([]cxEmailAccountDTO, 0, len(accounts))
	for i := range accounts {
		out = append(out, toCXEmailAccountDTO(&accounts[i]))
	}
	response.JSON(w, http.StatusOK, out)
}

func (h *CXHandler) CreateEmailAccount(w http.ResponseWriter, r *http.Request) {
	var in cxEmailAccountInput
	if !cxDecode(w, r, &in) {
		return
	}
	a, err := h.svc.CreateEmailAccount(r.Context(), cxTenant(r), &domain.EmailAccount{Name: in.Name, SMTPHost: in.SMTPHost, SMTPPort: in.SMTPPort, SMTPUsername: in.SMTPUsername, SMTPFrom: in.SMTPFrom, IMAPHost: in.IMAPHost, IMAPPort: in.IMAPPort, IMAPUsername: in.IMAPUsername, Active: in.Active})
	if err != nil {
		cxErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toCXEmailAccountDTO(a))
}

func (h *CXHandler) UpdateEmailAccount(w http.ResponseWriter, r *http.Request) {
	id, ok := cxID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in cxEmailAccountInput
	if !cxDecode(w, r, &in) {
		return
	}
	a, err := h.svc.UpdateEmailAccount(r.Context(), cxTenant(r), &domain.EmailAccount{ID: id, Name: in.Name, SMTPHost: in.SMTPHost, SMTPPort: in.SMTPPort, SMTPUsername: in.SMTPUsername, SMTPFrom: in.SMTPFrom, IMAPHost: in.IMAPHost, IMAPPort: in.IMAPPort, IMAPUsername: in.IMAPUsername, Active: in.Active})
	if err != nil {
		cxErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toCXEmailAccountDTO(a))
}

// IngestEmail is the authenticated ingestion point for an operator's
// IMAP poller (or a Mailpit test): it turns one inbound email into a
// message on the right conversation. It is staff-only in server.go,
// so unauthenticated senders cannot inject mail into the inbox.
func (h *CXHandler) IngestEmail(w http.ResponseWriter, r *http.Request) {
	var in cxEmailIngestInput
	if !cxDecode(w, r, &in) {
		return
	}
	c, err := h.svc.IngestEmail(r.Context(), cxTenant(r), in.AccountID, in.FromEmail, in.Subject, in.Body, in.ThreadID, in.ExternalMessageID)
	if err != nil {
		cxErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toCXConversationDTO(c))
}

func (h *CXHandler) ListWhatsAppAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.svc.ListWhatsAppAccounts(r.Context(), cxTenant(r))
	if err != nil {
		cxErr(w, err)
		return
	}
	out := make([]cxWhatsAppAccountDTO, 0, len(accounts))
	for i := range accounts {
		out = append(out, toCXWhatsAppAccountDTO(&accounts[i]))
	}
	response.JSON(w, http.StatusOK, out)
}

func (h *CXHandler) CreateWhatsAppAccount(w http.ResponseWriter, r *http.Request) {
	var in cxWhatsAppAccountInput
	if !cxDecode(w, r, &in) {
		return
	}
	a, err := h.svc.CreateWhatsAppAccount(r.Context(), cxTenant(r), &domain.WhatsAppAccount{Name: in.Name, BusinessID: in.BusinessID, PhoneNumberID: in.PhoneNumberID, DisplayName: in.DisplayName, Active: in.Active}, in.VerifyToken)
	if err != nil {
		cxErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toCXWhatsAppAccountDTO(a))
}

func (h *CXHandler) UpdateWhatsAppAccount(w http.ResponseWriter, r *http.Request) {
	id, ok := cxID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in cxWhatsAppAccountInput
	if !cxDecode(w, r, &in) {
		return
	}
	a, err := h.svc.UpdateWhatsAppAccount(r.Context(), cxTenant(r), &domain.WhatsAppAccount{ID: id, Name: in.Name, BusinessID: in.BusinessID, PhoneNumberID: in.PhoneNumberID, DisplayName: in.DisplayName, Active: in.Active})
	if err != nil {
		cxErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toCXWhatsAppAccountDTO(a))
}

// WhatsAppWebhookVerify answers Meta's verification challenge: the
// mode must be "subscribe" and the token must hash to a stored,
// active account's hash. The challenge is echoed back verbatim, as
// Meta requires; nothing else is returned.
func (h *CXHandler) WhatsAppWebhookVerify(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	_, err := h.svc.VerifyWhatsAppWebhook(r.Context(), cxTenant(r), q.Get("hub.mode"), q.Get("hub.verify_token"))
	if err != nil {
		cxErr(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	//nolint:gosec // G705: Meta webhook verification echoes hub.challenge as text/plain, never as HTML
	_, _ = w.Write([]byte(q.Get("hub.challenge")))
}

// WhatsAppWebhookIngest receives Meta message webhooks. It is
// authenticated by the X-Hub-Signature-256 HMAC over the raw body
// with the operator's app secret (CX_WHATSAPP_APP_SECRET): without a
// configured secret the endpoint refuses everything, so no message
// can be injected before activation. The service ingests each text
// message idempotently by its Meta message id; delivery statuses and
// non-text messages are acknowledged and ignored.
func (h *CXHandler) WhatsAppWebhookIngest(w http.ResponseWriter, r *http.Request) {
	secret := os.Getenv("CX_WHATSAPP_APP_SECRET")
	if secret == "" {
		response.Error(w, http.StatusServiceUnavailable, "whatsapp not configured")
		return
	}
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	signature := strings.TrimPrefix(r.Header.Get("X-Hub-Signature-256"), "sha256=")
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(raw)
	expected := hex.EncodeToString(mac.Sum(nil))
	if signature == "" || !hmac.Equal([]byte(signature), []byte(expected)) {
		response.Error(w, http.StatusForbidden, "invalid signature")
		return
	}
	var payload struct {
		Entry []struct {
			Changes []struct {
				Value struct {
					Messages []struct {
						From string `json:"from"`
						ID   string `json:"id"`
						Type string `json:"type"`
						Text struct {
							Body string `json:"body"`
						} `json:"text"`
					} `json:"messages"`
				} `json:"value"`
			} `json:"changes"`
		} `json:"entry"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	ingested := 0
	for _, entry := range payload.Entry {
		for _, change := range entry.Changes {
			for _, msg := range change.Value.Messages {
				if msg.Type != "text" || msg.Text.Body == "" {
					continue
				}
				threadID := "whatsapp:" + msg.From
				// The account id is not resolved from metadata yet
				// (see docs/cx-integraciones.md): the service does
				// not use it for routing, which is by tenant and
				// sender phone. Passing Nil keeps the call honest
				// about that gap instead of inventing an account.
				if _, err := h.svc.IngestWhatsApp(r.Context(), cxTenant(r), uuid.Nil, msg.From, msg.Text.Body, threadID, msg.ID); err != nil {
					if errors.Is(err, domain.ErrConflict) {
						continue // duplicate delivery, already ingested
					}
					cxErr(w, err)
					return
				}
				ingested++
			}
		}
	}
	response.JSON(w, http.StatusOK, map[string]any{"ingested": ingested})
}

// ── SLA policies, assignment rules, auto-replies ────────────────────

func (h *CXHandler) ListSLAPolicies(w http.ResponseWriter, r *http.Request) {
	policies, err := h.svc.ListSLAPolicies(r.Context(), cxTenant(r))
	if err != nil {
		cxErr(w, err)
		return
	}
	out := make([]cxSLAPolicyDTO, 0, len(policies))
	for i := range policies {
		out = append(out, toCXSLAPolicyDTO(&policies[i]))
	}
	response.JSON(w, http.StatusOK, out)
}

func (h *CXHandler) UpsertSLAPolicy(w http.ResponseWriter, r *http.Request) {
	var in cxSLAPolicyInput
	if !cxDecode(w, r, &in) {
		return
	}
	p, err := h.svc.UpsertSLAPolicy(r.Context(), cxTenant(r), &domain.SLAPolicy{Channel: in.Channel, Priority: in.Priority, ResponseMinutes: in.ResponseMinutes, ResolutionMinutes: in.ResolutionMinutes})
	if err != nil {
		cxErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toCXSLAPolicyDTO(p))
}

func (h *CXHandler) ListAssignmentRules(w http.ResponseWriter, r *http.Request) {
	rules, err := h.svc.ListAssignmentRules(r.Context(), cxTenant(r))
	if err != nil {
		cxErr(w, err)
		return
	}
	out := make([]cxAssignmentRuleDTO, 0, len(rules))
	for i := range rules {
		out = append(out, toCXAssignmentRuleDTO(&rules[i]))
	}
	response.JSON(w, http.StatusOK, out)
}

func (h *CXHandler) CreateAssignmentRule(w http.ResponseWriter, r *http.Request) {
	var in cxAssignmentRuleInput
	if !cxDecode(w, r, &in) {
		return
	}
	rule, err := h.svc.CreateAssignmentRule(r.Context(), cxTenant(r), &domain.AssignmentRule{Channel: in.Channel, Priority: in.Priority, AssigneeUserID: in.AssigneeUserID, Position: in.Position, Active: in.Active})
	if err != nil {
		cxErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toCXAssignmentRuleDTO(rule))
}

func (h *CXHandler) DeleteAssignmentRule(w http.ResponseWriter, r *http.Request) {
	id, ok := cxID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteAssignmentRule(r.Context(), cxTenant(r), id); err != nil {
		cxErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CXHandler) ListAutoReplies(w http.ResponseWriter, r *http.Request) {
	replies, err := h.svc.ListAutoReplies(r.Context(), cxTenant(r))
	if err != nil {
		cxErr(w, err)
		return
	}
	out := make([]cxAutoReplyDTO, 0, len(replies))
	for i := range replies {
		out = append(out, toCXAutoReplyDTO(&replies[i]))
	}
	response.JSON(w, http.StatusOK, out)
}

func (h *CXHandler) CreateAutoReply(w http.ResponseWriter, r *http.Request) {
	var in cxAutoReplyInput
	if !cxDecode(w, r, &in) {
		return
	}
	reply, err := h.svc.CreateAutoReply(r.Context(), cxTenant(r), &domain.AutoReply{TriggerKind: in.TriggerKind, Keyword: in.Keyword, Channel: in.Channel, ResponseText: in.ResponseText, Active: in.Active})
	if err != nil {
		cxErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toCXAutoReplyDTO(reply))
}

func (h *CXHandler) DeleteAutoReply(w http.ResponseWriter, r *http.Request) {
	id, ok := cxID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteAutoReply(r.Context(), cxTenant(r), id); err != nil {
		cxErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ── Surveys and reports ─────────────────────────────────────────────

func (h *CXHandler) CreateSurvey(w http.ResponseWriter, r *http.Request) {
	id, ok := cxID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in cxSurveyInput
	if !cxDecode(w, r, &in) {
		return
	}
	viewerID, role := cxCaller(r)
	s, err := h.svc.CreateSurvey(r.Context(), cxTenant(r), viewerID, role, id, &domain.Survey{Rating: in.Rating, Comment: in.Comment})
	if err != nil {
		cxErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toCXSurveyDTO(s))
}

func (h *CXHandler) ListSurveys(w http.ResponseWriter, r *http.Request) {
	viewerID, role := cxCaller(r)
	surveys, err := h.svc.ListSurveys(r.Context(), cxTenant(r), viewerID, role, cxOptID(r.URL.Query().Get("customer_id")))
	if err != nil {
		cxErr(w, err)
		return
	}
	out := make([]cxSurveyDTO, 0, len(surveys))
	for i := range surveys {
		out = append(out, toCXSurveyDTO(&surveys[i]))
	}
	response.JSON(w, http.StatusOK, out)
}

func (h *CXHandler) Report(w http.ResponseWriter, r *http.Request) {
	viewerID, role := cxCaller(r)
	report, err := h.svc.Report(r.Context(), cxTenant(r), viewerID, role)
	if err != nil {
		cxErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, report)
}
