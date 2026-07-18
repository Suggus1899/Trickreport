package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	appAnalytics "github.com/trickreport/backend/internal/application/analytics"
	appArticle "github.com/trickreport/backend/internal/application/article"
	appAuth "github.com/trickreport/backend/internal/application/auth"
	appAuto "github.com/trickreport/backend/internal/application/automation"
	appSLA "github.com/trickreport/backend/internal/application/sla"
	appTicket "github.com/trickreport/backend/internal/application/ticket"
	appUser "github.com/trickreport/backend/internal/application/user"
	domainArticle "github.com/trickreport/backend/internal/domain/article"
	domainTicket "github.com/trickreport/backend/internal/domain/ticket"
	"github.com/trickreport/backend/internal/interfaces/http/middleware"
	"github.com/trickreport/backend/internal/interfaces/http/response"
	"github.com/trickreport/backend/internal/interfaces/http/views"

	"github.com/a-h/templ"
)

// WebHandler serves HTML pages (server-side rendered with templ).
type WebHandler struct {
	auth       *appAuth.Service
	users      *appUser.Service
	tickets    *appTicket.UserService
	articles   *appArticle.Service
	sla        *appSLA.Service
	automation *appAuto.Service
	analytics  *appAnalytics.Service
}

func NewWebHandler(
	auth *appAuth.Service,
	users *appUser.Service,
	tickets *appTicket.UserService,
	articles *appArticle.Service,
	sla *appSLA.Service,
	automation *appAuto.Service,
	analytics *appAnalytics.Service,
) *WebHandler {
	return &WebHandler{
		auth: auth, users: users, tickets: tickets,
		articles: articles, sla: sla, automation: automation, analytics: analytics,
	}
}

// viewData builds common ViewData from the request context.
func (h *WebHandler) viewData(r *http.Request, title string) views.ViewData {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	d := views.ViewData{Title: title, Path: r.URL.Path}
	if claims != nil {
		d.UserRole = claims.Role
		d.TenantID = claims.TenantID
		// Fetch user name
		if u, err := h.users.Get(r.Context(), claims.UserID, claims.TenantID); err == nil {
			d.UserName = u.Name
		}
	}
	return d
}

// render wraps content in the Layout and writes it.
func (h *WebHandler) render(w http.ResponseWriter, r *http.Request, title string, content templ.Component) {
	data := views.LayoutData{
		Title:    title,
		UserName: "",
		UserRole: "",
		Path:     r.URL.Path,
		Content:  content,
	}
	claims, _ := middleware.ClaimsFromContext(r.Context())
	if claims != nil {
		data.UserRole = claims.Role
		if u, err := h.users.Get(r.Context(), claims.UserID, claims.TenantID); err == nil {
			data.UserName = u.Name
		}
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Layout(data).Render(r.Context(), w)
}

// ─── Auth pages ─────────────────────────────────────────────────────────────

func (h *WebHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	errorMsg := r.URL.Query().Get("error")
	views.Login(errorMsg).Render(r.Context(), w)
}

func (h *WebHandler) LoginSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/login?error=invalid+form", http.StatusSeeOther)
		return
	}

	result, err := h.auth.Login(r.Context(), appAuth.LoginInput{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	})
	if err != nil {
		http.Redirect(w, r, "/login?error=invalid+credentials", http.StatusSeeOther)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "trickreport_token", Value: result.Token, Path: "/",
		HttpOnly: true, Secure: result.SecureCookie,
		SameSite: http.SameSiteLaxMode, Expires: result.ExpiresAt,
	})

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *WebHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name: "trickreport_token", Value: "", Path: "/",
		MaxAge: -1, Expires: time.Unix(0, 0),
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// ─── Dashboard ──────────────────────────────────────────────────────────────

func (h *WebHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())
	summary, _ := h.analytics.GetSummary(r.Context(), tenantID)
	h.render(w, r, "Dashboard", views.Dashboard(views.SummaryData{
		TotalTickets: summary.TotalTickets, OpenTickets: summary.OpenTickets,
		ResolvedTickets: summary.ResolvedTickets, SLABreached: summary.SLABreached,
	}))
}

// ─── Tickets ────────────────────────────────────────────────────────────────

func (h *WebHandler) TicketList(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	tenantID := middleware.TenantFromContext(r.Context())

	tickets, err := h.tickets.List(r.Context(), tenantID, appTicket.Filter{}, claims.Role, claims.UserID)
	if err != nil {
		h.render(w, r, "Tickets", views.TicketList(nil))
		return
	}

	dtos := make([]views.TicketData, len(tickets))
	for i, t := range tickets {
		dtos[i] = ticketToView(t)
	}
	h.render(w, r, "Tickets", views.TicketList(dtos))
}

func (h *WebHandler) TicketDetailPage(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	tenantID := middleware.TenantFromContext(r.Context())
	ticketID := chi.URLParam(r, "id")

	t, err := h.tickets.Get(r.Context(), ticketID, tenantID, claims.Role, claims.UserID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	comments, _ := h.tickets.ListComments(r.Context(), ticketID, tenantID, claims.Role, claims.UserID)

	var history []views.HistoryData
	if claims.Role == "admin" || claims.Role == "agent" {
		entries, _ := h.tickets.ListHistory(r.Context(), ticketID, tenantID, claims.Role, claims.UserID)
		history = make([]views.HistoryData, len(entries))
		for i, e := range entries {
			history[i] = views.HistoryData{
				ID: e.ID, UserID: e.UserID, UserName: e.UserName,
				Field: e.Field, OldValue: e.OldValue, NewValue: e.NewValue,
				CreatedAt: e.CreatedAt,
			}
		}
	}

	commentDtos := make([]views.CommentData, len(comments))
	for i, c := range comments {
		commentDtos[i] = views.CommentData{
			ID: c.ID, TicketID: c.TicketID, UserID: c.UserID,
			UserName: c.UserName, Content: c.Content,
			IsInternal: c.IsInternal, CreatedAt: c.CreatedAt,
		}
	}

	canComment := true
	h.render(w, r, t.Title, views.TicketDetail(ticketToView(*t), commentDtos, history, canComment))
}

func (h *WebHandler) TicketFormPage(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, "New Ticket", views.TicketForm())
}

func (h *WebHandler) TicketCreate(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	tenantID := middleware.TenantFromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/tickets/new", http.StatusSeeOther)
		return
	}

	t, err := h.tickets.Create(r.Context(), appTicket.CreateInput{
		TenantID:    tenantID,
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Priority:    r.FormValue("priority"),
		Category:    r.FormValue("category"),
		CreatedBy:   claims.UserID,
	})
	if err != nil {
		http.Redirect(w, r, "/tickets/new", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/tickets/"+t.ID, http.StatusSeeOther)
}

func (h *WebHandler) TicketAddComment(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	tenantID := middleware.TenantFromContext(r.Context())
	ticketID := chi.URLParam(r, "id")

	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/tickets/"+ticketID, http.StatusSeeOther)
		return
	}

	_, err := h.tickets.AddComment(r.Context(), appTicket.AddCommentInput{
		TicketID: ticketID, TenantID: tenantID,
		UserID: claims.UserID, Role: claims.Role,
		Content: r.FormValue("content"),
	})
	if err != nil {
		http.Redirect(w, r, "/tickets/"+ticketID, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/tickets/"+ticketID, http.StatusSeeOther)
}

// ─── Articles ───────────────────────────────────────────────────────────────

func (h *WebHandler) ArticleList(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	tenantID := middleware.TenantFromContext(r.Context())

	articles, err := h.articles.List(r.Context(), tenantID, appArticle.Filter{}, claims.Role)
	if err != nil {
		h.render(w, r, "Knowledge Base", views.ArticleList(nil))
		return
	}

	dtos := make([]views.ArticleData, len(articles))
	for i, a := range articles {
		dtos[i] = articleToView(a)
	}
	h.render(w, r, "Knowledge Base", views.ArticleList(dtos))
}

func (h *WebHandler) ArticleDetailPage(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	tenantID := middleware.TenantFromContext(r.Context())
	id := chi.URLParam(r, "id")

	a, err := h.articles.Get(r.Context(), id, tenantID, claims.Role)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	h.render(w, r, a.Title, views.ArticleDetail(articleToView(*a)))
}

func (h *WebHandler) ArticleFormPage(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, "New Article", views.ArticleForm(nil))
}

func (h *WebHandler) ArticleCreate(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	tenantID := middleware.TenantFromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/articles/new", http.StatusSeeOther)
		return
	}

	tags := parseTags(r.FormValue("tags"))
	published := r.FormValue("published") == "true"

	a, err := h.articles.Create(r.Context(), appArticle.CreateInput{
		TenantID: tenantID, Title: r.FormValue("title"), Content: r.FormValue("content"),
		Category: r.FormValue("category"), Tags: tags, Published: published,
		CreatedBy: claims.UserID,
	})
	if err != nil {
		http.Redirect(w, r, "/articles/new", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/articles/"+a.ID, http.StatusSeeOther)
}

// ─── Admin: Users ───────────────────────────────────────────────────────────

func (h *WebHandler) UserList(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())

	users, err := h.users.List(r.Context(), tenantID)
	if err != nil {
		h.render(w, r, "Users", views.UserList(nil))
		return
	}

	dtos := make([]views.UserData, len(users))
	for i, u := range users {
		dtos[i] = views.UserData{
			ID: u.ID, Name: u.Name, Email: u.Email, Role: string(u.Role),
			AvatarURL: u.AvatarURL, Active: u.Active, CreatedAt: u.CreatedAt,
		}
	}
	h.render(w, r, "Users", views.UserList(dtos))
}

func (h *WebHandler) UserFormPage(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, "New User", views.UserForm(nil))
}

func (h *WebHandler) UserCreate(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/admin/users/new", http.StatusSeeOther)
		return
	}

	_, err := h.users.Create(r.Context(), appUser.CreateInput{
		TenantID: tenantID, Name: r.FormValue("name"), Email: r.FormValue("email"),
		Role: r.FormValue("role"), Password: r.FormValue("password"),
	})
	if err != nil {
		http.Redirect(w, r, "/admin/users/new", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

// ─── Admin: SLA ─────────────────────────────────────────────────────────────

func (h *WebHandler) SLAList(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())

	policies, err := h.sla.List(r.Context(), tenantID)
	if err != nil {
		h.render(w, r, "SLA Policies", views.SLAList(nil))
		return
	}

	dtos := make([]views.SLAPolicyData, len(policies))
	for i, p := range policies {
		dtos[i] = views.SLAPolicyData{
			Priority: p.Priority, ResponseTimeMinutes: p.ResponseTimeMinutes,
			ResolutionTimeMinutes: p.ResolutionTimeMinutes, EscalationMinutes: p.EscalationMinutes,
		}
	}
	h.render(w, r, "SLA Policies", views.SLAList(dtos))
}

func (h *WebHandler) SLAFormPage(w http.ResponseWriter, r *http.Request) {
	priority := chi.URLParam(r, "priority")
	h.render(w, r, "SLA Policy", views.SLAForm(priority, nil))
}

func (h *WebHandler) SLAUpsert(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())
	priority := chi.URLParam(r, "priority")

	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/admin/sla/"+priority+"/edit", http.StatusSeeOther)
		return
	}

	_, err := h.sla.Upsert(r.Context(), appSLA.UpsertInput{
		TenantID: tenantID, Priority: priority,
		ResponseTimeMinutes:   formInt(r, "response_time_minutes"),
		ResolutionTimeMinutes: formInt(r, "resolution_time_minutes"),
		EscalationMinutes:     formInt(r, "escalation_minutes"),
	})
	if err != nil {
		http.Redirect(w, r, "/admin/sla/"+priority+"/edit", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin/sla", http.StatusSeeOther)
}

// ─── Admin: Automations ─────────────────────────────────────────────────────

func (h *WebHandler) AutomationList(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())

	rules, err := h.automation.List(r.Context(), tenantID)
	if err != nil {
		h.render(w, r, "Automation Rules", views.AutomationList(nil))
		return
	}

	dtos := make([]views.AutomationData, len(rules))
	for i, rule := range rules {
		dtos[i] = views.AutomationData{
			ID: rule.ID, Name: rule.Name, Description: rule.Description,
			TriggerType: rule.TriggerType, IsActive: rule.IsActive, CreatedAt: rule.CreatedAt,
		}
	}
	h.render(w, r, "Automation Rules", views.AutomationList(dtos))
}

func (h *WebHandler) AutomationFormPage(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, "New Automation Rule", views.AutomationForm(nil))
}

func (h *WebHandler) AutomationCreate(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/admin/automations/new", http.StatusSeeOther)
		return
	}

	conditions := parseJSONMap(r.FormValue("conditions"))
	actions := parseJSONArray(r.FormValue("actions"))

	_, err := h.automation.Create(r.Context(), appAuto.CreateInput{
		TenantID: tenantID, Name: r.FormValue("name"), Description: r.FormValue("description"),
		TriggerType: r.FormValue("trigger_type"), Conditions: conditions, Actions: actions,
		IsActive: r.FormValue("is_active") == "true",
	})
	if err != nil {
		http.Redirect(w, r, "/admin/automations/new", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin/automations", http.StatusSeeOther)
}

// ─── Admin: Analytics ───────────────────────────────────────────────────────

func (h *WebHandler) AnalyticsPage(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())

	summary, _ := h.analytics.GetSummary(r.Context(), tenantID)
	volume, _ := h.analytics.GetVolume(r.Context(), tenantID)
	statusDist, _ := h.analytics.GetStatusDistribution(r.Context(), tenantID)
	resolution, _ := h.analytics.GetResolutionTime(r.Context(), tenantID)

	volDtos := make([]views.VolumePointData, len(volume))
	for i, v := range volume {
		volDtos[i] = views.VolumePointData{Date: v.Date, Count: v.Count}
	}

	statusDtos := make([]views.StatusDistData, len(statusDist))
	for i, s := range statusDist {
		statusDtos[i] = views.StatusDistData{Status: s.Status, Count: s.Count}
	}

	resDtos := make([]views.ResolutionData, len(resolution))
	for i, res := range resolution {
		resDtos[i] = views.ResolutionData{Priority: res.Priority, AvgHours: res.AvgHours}
	}

	h.render(w, r, "Analytics", views.Analytics(
		views.SummaryData{
			TotalTickets: summary.TotalTickets, OpenTickets: summary.OpenTickets,
			ResolvedTickets: summary.ResolvedTickets, SLABreached: summary.SLABreached,
		},
		volDtos, statusDtos, resDtos,
	))
}

// ─── Helpers ────────────────────────────────────────────────────────────────

func ticketToView(t domainTicket.Ticket) views.TicketData {
	return views.TicketData{
		ID: t.ID, Title: t.Title, Description: t.Description,
		Status: string(t.Status), Priority: string(t.Priority), Category: t.Category,
		CreatedBy: t.CreatedBy, CreatorName: t.CreatorName,
		AssignedTo: t.AssignedTo, AssigneeName: t.AssigneeName,
		SLADeadline: t.SLADeadline, SLABreached: t.SLABreached,
		CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt,
	}
}

func articleToView(a domainArticle.Article) views.ArticleData {
	return views.ArticleData{
		ID: a.ID, Title: a.Title, Content: a.Content, Category: a.Category,
		Tags: a.Tags, Published: a.Published, CreatedBy: a.CreatedBy,
		AuthorName: a.AuthorName, CreatedAt: a.CreatedAt, UpdatedAt: a.UpdatedAt,
	}
}

func parseTags(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

func formInt(r *http.Request, field string) int {
	v := r.FormValue(field)
	n := 0
	for _, c := range v {
		if c < '0' || c > '9' {
			return 0
		}
		n = n*10 + int(c-'0')
	}
	return n
}

func parseJSONMap(s string) map[string]any {
	if s == "" {
		return nil
	}
	var m map[string]any
	_ = json.Unmarshal([]byte(s), &m)
	return m
}

func parseJSONArray(s string) []any {
	if s == "" {
		return nil
	}
	var a []any
	_ = json.Unmarshal([]byte(s), &a)
	return a
}

// Ensure imports are used
var (
	_ = errors.Is
	_ = context.Background
	_ = response.JSON
)
