package mailpen

import (
	"fmt"
)

// The following structs represent the data needed to render various components in an email templates.

// commonEmailData adds common data to the email data map.
func commonTemplateData(cfg *Config, data map[string]any) map[string]any {
	data["FooterData"] = FooterData{
		CompanyName:   cfg.CompanyName,
		SupportEmail:  cfg.SupportEmail,
		CopyrightText: fmt.Sprintf("© 2024 %s. All rights reserved.", cfg.CompanyName),
		AddressLine1:  cfg.CompanyAddress1,
	}

	return data
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

// TableData represents the data needed to render a table
type TableData struct {
	Headers []TableHeader
	Rows    []TableRow
}

// TwoColumnRow represents a row in a two-column layout
type TwoColumnRow struct {
	Label string
	Value string
}

// TwoColumnData represents the data needed to render a two-column layout
type TwoColumnData struct {
	Rows []TwoColumnRow
}

// FooterData represents the data needed to render a footer
type FooterData struct {
	CompanyName   string
	SupportEmail  string
	CopyrightText string // e.g., "© 2024 Crystal Springs Foundation. All rights reserved."
	AddressLine1  string // e.g., "1234 Business Street, Suite 500"
	AddressLine2  string // e.g., "San Francisco, CA 94111"
}

// NotificationButton represents the type of button to render in a notification box
type NotificationButton struct {
	BgColor     string
	BorderColor string
	TextColor   string
	Text        string
	URL         string
}

// NotificationBoxData represents the data needed to render a notification box
type NotificationBoxData struct {
	BgColor     string // e.g., "#FFF3CD" for warning
	BorderColor string // e.g., "#FFA500" for warning
	Icon        string // Optional icon URL
	IconAlt     string
	Title       string
	TitleColor  string
	Message     string
	TextColor   string
	Button      *NotificationButton
}

// Card represents a card in a card grid
type Card struct {
	ImageURL    string
	ImageAlt    string
	Title       string
	Description string
	LinkURL     string
	LinkText    string
}

// CardGridData represents the data needed to render a card grid
type CardGridData struct {
	Cards []Card
}
