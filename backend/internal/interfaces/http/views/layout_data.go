package views

import "github.com/a-h/templ"

// LayoutData holds the data for the base HTML layout.
type LayoutData struct {
	Title    string
	UserName string
	UserRole string
	Path     string
	Content  templ.Component
}
