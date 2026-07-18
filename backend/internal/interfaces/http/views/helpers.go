package views

import "strconv"

// NavItem represents a navigation menu entry.
type NavItem struct {
	Label  string
	Path   string
	Icon   string
	Roles  []string // empty = all roles
	Active bool
}

// navItems returns the navigation items filtered by role.
func navItems(currentPath, role string) []NavItem {
	all := []NavItem{
		{Label: "Dashboard", Path: "/dashboard", Icon: "M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"},
		{Label: "Tickets", Path: "/tickets", Icon: "M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"},
		{Label: "Knowledge Base", Path: "/articles", Icon: "M12 6.253v13m0-13C10.015 5.422 7.252 5 4.752 5c-1.5 0-2.748.423-3.752 1.253v13C2.254 18.423 3.502 18 5.002 18c2.5 0 4.763.422 6.748 1.253m0-13C13.985 5.422 16.748 5 19.248 5c1.5 0 2.748.423 3.752 1.253v13C22.246 18.423 20.998 18 19.498 18c-2.5 0-4.763.422-6.748 1.253"},
		{Label: "Users", Path: "/admin/users", Icon: "M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z", Roles: []string{"admin"}},
		{Label: "SLA Policies", Path: "/admin/sla", Icon: "M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z", Roles: []string{"admin"}},
		{Label: "Automations", Path: "/admin/automations", Icon: "M13 10V3L4 14h7v7l9-11h-7z", Roles: []string{"admin"}},
		{Label: "Analytics", Path: "/admin/analytics", Icon: "M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z", Roles: []string{"admin"}},
	}

	var result []NavItem
	for _, item := range all {
		visible := len(item.Roles) == 0
		for _, r := range item.Roles {
			if r == role {
				visible = true
				break
			}
		}
		if visible {
			item.Active = currentPath == item.Path || (item.Path != "/dashboard" && len(currentPath) > len(item.Path) && currentPath[:len(item.Path)] == item.Path)
			result = append(result, item)
		}
	}
	return result
}

// statusBadgeClass returns the CSS class for a ticket status badge.
func statusBadgeClass(status string) string {
	switch status {
	case "open":
		return "badge badge-blue"
	case "in_progress":
		return "badge badge-yellow"
	case "waiting_client":
		return "badge badge-purple"
	case "resolved":
		return "badge badge-green"
	case "closed":
		return "badge badge-gray"
	default:
		return "badge badge-gray"
	}
}

// priorityBadgeClass returns the CSS class for a ticket priority badge.
func priorityBadgeClass(priority string) string {
	switch priority {
	case "low":
		return "badge badge-gray"
	case "medium":
		return "badge badge-blue"
	case "high":
		return "badge badge-yellow"
	case "critical":
		return "badge badge-red"
	default:
		return "badge badge-gray"
	}
}

// formatTime formats a time.Time as "2006-01-02 15:04".
func formatTime(t any) string {
	type timeLike interface {
		Format(layout string) string
	}
	if tl, ok := t.(timeLike); ok {
		return tl.Format("2006-01-02 15:04")
	}
	return ""
}

// formatTimePtr formats a *time.Time, returning "—" if nil.
func formatTimePtr(t any) string {
	if t == nil {
		return "—"
	}
	return formatTime(t)
}

// ptrStr returns the string value of a *string, or "—" if nil.
func ptrStr(s *string) string {
	if s == nil {
		return "—"
	}
	return *s
}

// intStr converts an int to string.
func intStr(i int) string {
	return strconv.Itoa(i)
}

// floatStr formats a float64 with 1 decimal place.
func floatStr(f float64) string {
	return strconv.FormatFloat(f, 'f', 1, 64)
}

// ─── Article helpers ────────────────────────────────────────────────────────

func articleExcerpt(content string, maxLen int) string {
	if len(content) <= maxLen {
		return content
	}
	return content[:maxLen] + "..."
}

func articleTitle(a *ArticleData) string {
	if a == nil {
		return ""
	}
	return a.Title
}

func articleContent(a *ArticleData) string {
	if a == nil {
		return ""
	}
	return a.Content
}

func articleCategory(a *ArticleData) string {
	if a == nil {
		return "general"
	}
	return a.Category
}

func articleTags(a *ArticleData) string {
	if a == nil || len(a.Tags) == 0 {
		return ""
	}
	result := ""
	for i, t := range a.Tags {
		if i > 0 {
			result += ", "
		}
		result += t
	}
	return result
}

func articlePublished(a *ArticleData) bool {
	if a == nil {
		return false
	}
	return a.Published
}

// ─── User helpers ───────────────────────────────────────────────────────────

func userField(u *UserData, field string) string {
	if u == nil {
		return ""
	}
	switch field {
	case "name":
		return u.Name
	case "email":
		return u.Email
	}
	return ""
}

func userRoleIs(u *UserData, role string) bool {
	if u == nil {
		return role == "end_user"
	}
	return u.Role == role
}

// ─── SLA helpers ────────────────────────────────────────────────────────────

func slaField(p *SLAPolicyData, field string) string {
	if p == nil {
		return "0"
	}
	switch field {
	case "response":
		return intStr(p.ResponseTimeMinutes)
	case "resolution":
		return intStr(p.ResolutionTimeMinutes)
	case "escalation":
		return intStr(p.EscalationMinutes)
	}
	return "0"
}

// ─── Automation helpers ─────────────────────────────────────────────────────

func autoField(r *AutomationData, field string) string {
	if r == nil {
		return ""
	}
	switch field {
	case "name":
		return r.Name
	case "description":
		return r.Description
	case "conditions":
		return "{}"
	case "actions":
		return "[]"
	}
	return ""
}

func autoTriggerIs(r *AutomationData, trigger string) bool {
	if r == nil {
		return trigger == "ticket_created"
	}
	return r.TriggerType == trigger
}

func autoActive(r *AutomationData) bool {
	if r == nil {
		return true
	}
	return r.IsActive
}
