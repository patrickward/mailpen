package templates

import "html/template"

// Component represents a reusable email component
type Component interface {
	// Template returns the templates name
	Template() string

	// Data returns the data for the templates
	Data() any
}

// Layout represents an email layout
type Layout interface {
	// Template returns the templates name
	Template() string

	// Data returns the data for the templates
	Data() any

	// Components returns the components for the layout
	Components() []Component
}

// BaseLayout is the base layout for all layouts
type BaseLayout struct {
	Subject string
	Content template.HTML
	//Footer FooterData
}

// TableHeader represents a header in a table
type TableHeader struct {
	Text  string
	Width string
}

// TableCell represents a cell in a table
type TableCell struct {
	Text  string
	Width string
}

// TableRow represents a row in a table
type TableRow struct {
	Cells []TableCell
}

// Table implements Component
type Table struct {
	Headers []TableHeader
	Rows    []TableRow
}

func (t *Table) TemplateName() string {
	return "@table"
}

func (t *Table) Data() any {
	return t
}

// TwoColumnRow represents a row in a two-column layout
type TwoColumnRow struct {
	Label string
	Value string
}

// TwoColumn implements Component
type TwoColumn struct {
	Rows []TwoColumnRow
}

func (t *TwoColumn) TemplateName() string {
	return "@two-column"
}

func (t *TwoColumn) Data() any {
	return t
}

// Button implements Component
type Button struct {
	Text        string
	URL         string
	Style       string // default, primary, danger
	BgColor     string
	BorderColor string
	TextColor   string
}

func (b *Button) TemplateName() string {
	return "@button"
}

func (b *Button) Data() any {
	return b
}

// NotificationBox implements Component
type NotificationBox struct {
	BgColor     string
	BorderColor string
	Icon        string
	IconAlt     string
	Title       string
	TitleColor  string
	Message     string
	TextColor   string
	Button      *Button
}

func (n *NotificationBox) TemplateName() string {
	return "@notification-box"
}

func (n *NotificationBox) Data() any {
	return n
}

// Helper function to merge templates data with layout data

func MergeData(baseData, componentData map[string]any) map[string]any {
	result := make(map[string]any)

	// Copy base data
	for k, v := range baseData {
		result[k] = v
	}

	// Add/override with component data
	for k, v := range componentData {
		result[k] = v
	}

	return result
}
