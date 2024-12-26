package mailpen

import (
	"html/template"
)

// Config holds the mailpen configuration
type Config struct {
	// From address
	From    string // From address
	ReplyTo string // Reply-to address

	// Company/Branding
	BaseURL         string // Base URL of the website
	CompanyAddress1 string // The first line of the company address (usually the street address)
	CompanyAddress2 string // The second line of the company address (usually the city, state, and ZIP code)
	CompanyName     string // Company name
	LogoURL         string // URL to the company logo
	SupportEmail    string // Support email address
	SupportPhone    string // Support phone number
	WebsiteName     string // Name of the website
	WebsiteURL      string // URL to the company website.

	// HTML processor for processing HTML content
	HTMLProcessor HTMLProcessor // HTML processor for processing HTML content

	// Links
	SiteLinks        map[string]string // Site links
	SocialMediaLinks map[string]string // Social media links

	// Template configuration
	Extensions []string         // Extensions for the templates. Defaults to []string{".html"}.
	FuncMap    template.FuncMap // Additional template functions to add to the template engine. These will be merged with the default functions.
	Sources    []TemplateSource // Template sources
}
