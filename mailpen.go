package mailpen

import (
	"context"
	"errors"
	"fmt"

	gomail "github.com/wneessen/go-mail"

	"github.com/patrickward/mailpen/templates"
)

var (
	ErrNoContent = errors.New("email must have either plain text or HTML body")
	ErrNoSubject = errors.New("email must have a subject")
)

// SMTPClient defines the interface for an SMTP client, mainly used for testing
type SMTPClient interface {
	DialAndSend(messages ...*gomail.Msg) error
}

// HTMLProcessor defines the interface for processing HTML content
type HTMLProcessor interface {
	Process(html string) (string, error)
}

// StringList is an alias for a slice of strings
type StringList = []string

// Option is a functional option for configuring a Mailpen instance
type Option func(mailpen *Mailpen) error

// Mailpen handles email sending operations
type Mailpen struct {
	config        *Config
	provider      Provider
	templateMgr   *Manager
	components    map[string]templates.Component
	htmlProcessor HTMLProcessor
}

// New creates a new Mailpen instance using the provided configuration and the default SMTP client
func New(provider Provider, config *Config, opts ...Option) (*Mailpen, error) {
	if provider == nil {
		return nil, errors.New("provider is required")
	}

	if config == nil {
		return nil, errors.New("config is required")
	}

	tmOpts := &ManagerConfig{
		FuncMap:   config.FuncMap,
		Processor: config.HTMLProcessor,
		Sources:   config.Sources,
	}

	tm, err := NewManager(tmOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create templates manager: %w", err)
	}

	mp := &Mailpen{
		config:      config,
		provider:    provider,
		templateMgr: tm,
		components:  make(map[string]templates.Component),
	}

	// Register built-in components
	if err := mp.registerBuiltinComponents(); err != nil {
		return nil, fmt.Errorf("failed to register built-in components: %w", err)
	}

	// Apply additional template sources
	if err := mp.addTemplateSources(config.Sources); err != nil {
		return nil, fmt.Errorf("failed to add template sources: %w", err)
	}

	// Apply options
	for _, opt := range opts {
		if err := opt(mp); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return mp, nil
}

// addTemplateSource adds a new template source to the templates manager
func (m *Mailpen) addTemplateSource(source TemplateSource) error {
	return m.templateMgr.AddSource(source)
}

// addTemplateSources adds additional template sources
func (m *Mailpen) addTemplateSources(sources []TemplateSource) error {
	for _, source := range sources {
		if err := m.addTemplateSource(source); err != nil {
			return fmt.Errorf("failed to add template source: %w", err)
		}
	}
	return nil
}

// RegisterComponent registers a new component
func (m *Mailpen) RegisterComponent(name string, component templates.Component) error {
	if _, exists := m.components[name]; exists {
		return fmt.Errorf("component %s already exists", name)
	}
	m.components[name] = component
	return nil
}

// GetComponent returns a registered component by name
func (m *Mailpen) GetComponent(name string) (templates.Component, bool) {
	component, exists := m.components[name]
	return component, exists
}

// Config returns the mailpen configuration
func (m *Mailpen) Config() *Config {
	return m.config
}

// Send sends an email using the provided templates and data
func (m *Mailpen) Send(ctx context.Context, msg *Message) error {
	if err := m.processTemplates(msg); err != nil {
		return fmt.Errorf("failed to process templates: %w", err)
	}

	if msg.From == "" {
		msg.From = m.config.From
	}

	// Send via provider
	return m.provider.Send(ctx, msg)
}

// NewTemplateData creates a new templates data map with default values
func (m *Mailpen) NewTemplateData() TemplateData {
	return NewTemplateData(m.config)
}

func (m *Mailpen) processTemplates(msg *Message) error {
	if msg.Template == "" {
		return nil
	}

	data := m.prepareTemplateData(msg.Data)

	rendered, err := m.templateMgr.RenderEmail(msg.Template, data, msg.Layout)
	if err != nil {
		return fmt.Errorf("failed to render email: %w", err)
	}

	if rendered.Subject != "" {
		msg.Subject = rendered.Subject
	}

	if rendered.Text != "" {
		msg.TextBody = rendered.Text
	}

	if rendered.HTML != "" {
		msg.HTMLBody = rendered.HTML
	}

	return nil
}

func (m *Mailpen) prepareTemplateData(data map[string]any) TemplateData {
	// Merge data with default values
	data = mergeData(m.NewTemplateData(), data)

	// Add global data
	data["Config"] = m.config

	return data
}

// registerBuiltinComponents registers the default components
func (m *Mailpen) registerBuiltinComponents() error {
	//components := []struct {
	//	name string
	//	comp templates.Component
	//}{
	//	{"button", &templates.Button{}},
	//	{"table", &templates.Table{}},
	//	{"two-column", &templates.TwoColumn{}},
	//	{"notification", &templates.NotificationBox{}},
	//	// Add more built-in components here
	//}
	//
	//for _, c := range components {
	//	if err := m.RegisterComponent(c.name, c.comp); err != nil {
	//		return err
	//	}
	//}

	return nil
}

// mergeData merges two data maps
func mergeData(base, overlay map[string]any) map[string]any {
	result := make(map[string]any)

	for k, v := range base {
		result[k] = v
	}

	for k, v := range overlay {
		result[k] = v
	}

	return result
}
