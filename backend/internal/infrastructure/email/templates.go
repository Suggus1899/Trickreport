package email

import (
	"bytes"
	"fmt"
	"html/template"
)

// TemplateData holds the common fields used by email templates.
type TemplateData struct {
	Title       string
	Description string
	Priority    string
	Status      string
	TicketID    string
	TenantName  string
	UserName    string
	Deadline    string
	ResetLink   string
}

var templates *template.Template

func init() {
	templates = template.Must(template.New("ticket_created").Parse(ticketCreatedTmpl))
	templates = template.Must(templates.New("ticket_updated").Parse(ticketUpdatedTmpl))
	templates = template.Must(templates.New("sla_breach").Parse(slaBreachTmpl))
	templates = template.Must(templates.New("password_reset").Parse(passwordResetTmpl))
}

// RenderTemplate renders the named HTML email template with the given data.
func RenderTemplate(name string, data TemplateData) (string, error) {
	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, name, data); err != nil {
		return "", fmt.Errorf("failed to render email template %q: %w", name, err)
	}
	return buf.String(), nil
}

const emailBase = `<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><title>{{.Title}}</title></head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
<h2 style="color: #333;">{{.Title}}</h2>
<div style="background: #f9f9f9; padding: 15px; border-radius: 5px;">%s</div>
<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
<p style="color: #999; font-size: 12px;">This is an automated message from Trickreport.</p>
</body>
</html>`

const ticketCreatedTmpl = `<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><title>New Ticket Created</title></head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
<h2 style="color: #333;">New Ticket Created</h2>
<div style="background: #f9f9f9; padding: 15px; border-radius: 5px;">
<p><strong>Ticket ID:</strong> {{.TicketID}}</p>
<p><strong>Title:</strong> {{.Title}}</p>
<p><strong>Priority:</strong> {{.Priority}}</p>
<p><strong>Description:</strong></p>
<p>{{.Description}}</p>
</div>
<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
<p style="color: #999; font-size: 12px;">This is an automated message from Trickreport.</p>
</body>
</html>`

const ticketUpdatedTmpl = `<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><title>Ticket Updated</title></head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
<h2 style="color: #333;">Ticket Updated</h2>
<div style="background: #f9f9f9; padding: 15px; border-radius: 5px;">
<p><strong>Ticket ID:</strong> {{.TicketID}}</p>
<p><strong>Title:</strong> {{.Title}}</p>
<p><strong>New Status:</strong> {{.Status}}</p>
</div>
<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
<p style="color: #999; font-size: 12px;">This is an automated message from Trickreport.</p>
</body>
</html>`

const slaBreachTmpl = `<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><title>SLA Breach Alert</title></head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
<h2 style="color: #c0392b;">SLA Breach Alert</h2>
<div style="background: #fdedee; padding: 15px; border-radius: 5px;">
<p><strong>Ticket ID:</strong> {{.TicketID}}</p>
<p><strong>Title:</strong> {{.Title}}</p>
<p><strong>Priority:</strong> {{.Priority}}</p>
<p><strong>Deadline:</strong> {{.Deadline}}</p>
<p>The resolution time SLA for this ticket has been breached.</p>
</div>
<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
<p style="color: #999; font-size: 12px;">This is an automated message from Trickreport.</p>
</body>
</html>`

const passwordResetTmpl = `<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><title>Password Reset</title></head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
<h2 style="color: #333;">Password Reset Request</h2>
<div style="background: #f9f9f9; padding: 15px; border-radius: 5px;">
<p>Hello {{.UserName}},</p>
<p>You requested a password reset. Click the link below to choose a new password:</p>
<p><a href="{{.ResetLink}}" style="color: #2980b9;">Reset Password</a></p>
<p>If you did not request this reset, you can safely ignore this email.</p>
</div>
<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
<p style="color: #999; font-size: 12px;">This is an automated message from Trickreport.</p>
</body>
</html>`

// Ensure the base template constant is referenced to avoid unused warnings.
var _ = emailBase
