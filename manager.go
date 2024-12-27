package mailpen

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"path"
	"strings"
	"sync"

	"github.com/patrickward/mailpen/templates"
)

const (
	LayoutsDir    = "layouts"
	PartialsDir   = "partials"
	ComponentsDir = "components"
	EmailsDir     = "emails"
)

// TemplateSource represents a source of templates
type TemplateSource struct {
	Name string // Name of the template source
	FS   fs.FS  // File system for the templates
}

// TemplateFormat represents the format of a template
type TemplateFormat string

const (
	FormatText TemplateFormat = "text"
	FormatHTML TemplateFormat = "html"
)

// Manager handles templates loading, caching, and rendering
type Manager struct {
	funcMap       template.FuncMap
	processor     HTMLProcessor
	defaultLayout string
	sources       []TemplateSource
	theme         map[string]any
	baseTemplates map[TemplateFormat]*template.Template
	emailCache    map[string]*template.Template
	mu            sync.RWMutex
}

// ManagerConfig configures the templates manager
type ManagerConfig struct {
	FuncMap       template.FuncMap
	Processor     HTMLProcessor
	Sources       []TemplateSource
	Theme         map[string]any
	DefaultLayout string
}

// DefaultProcessor provides a pass-through implementation
type DefaultProcessor struct{}

// Process provides a pass-through implementation for HTMLProcessor
func (p *DefaultProcessor) Process(html string) (string, error) {
	return html, nil
}

// NewManager creates a new templates manager
func NewManager(config *ManagerConfig) (*Manager, error) {
	if config == nil {
		config = &ManagerConfig{}
	}

	if config.Processor == nil {
		config.Processor = &DefaultProcessor{}
	}

	if config.DefaultLayout == "" {
		config.DefaultLayout = "base"
	}

	if config.Theme == nil {
		config.Theme = DefaultTheme()
	}

	m := &Manager{
		processor:     config.Processor,
		defaultLayout: config.DefaultLayout,
		sources:       make([]TemplateSource, 0),
		baseTemplates: make(map[TemplateFormat]*template.Template),
		emailCache:    make(map[string]*template.Template),
		theme:         config.Theme,
	}

	// Merge function maps
	m.funcMap = MergeFuncMaps(DefaultFuncMap(), m.funcMap, m.themeFuncs())

	// Initialize base template sets
	m.baseTemplates[FormatText] = template.New("text-base").Funcs(m.funcMap)
	m.baseTemplates[FormatHTML] = template.New("html-base").Funcs(m.funcMap)

	// Add the built-in templates as a source
	if err := m.AddSource(TemplateSource{
		Name: "built-in",
		FS:   templates.FS,
	}); err != nil {
		return nil, fmt.Errorf("failed to add built-in templates: %w", err)
	}

	// Add initial sources if provided
	for _, source := range config.Sources {
		if err := m.AddSource(source); err != nil {
			return nil, fmt.Errorf("failed to add source %q: %w", source.Name, err)
		}
	}

	return m, nil
}

// formatFromFile determines the template format from filename
func formatFromFile(filename string) TemplateFormat {
	ext := path.Ext(filename)
	switch ext {
	case ".html":
		return FormatHTML
	case ".txt":
		return FormatText
	default:
		return ""
	}
}

// loadBaseTemplates loads layouts, components, and partials
func (m *Manager) loadBaseTemplates() error {
	// Reset base templates
	m.baseTemplates[FormatText] = template.New("text-base").Funcs(m.funcMap)
	m.baseTemplates[FormatHTML] = template.New("html-base").Funcs(m.funcMap)

	// Load from each source in order
	for _, source := range m.sources {
		// Load layouts
		if err := m.loadDirectory(source, LayoutsDir); err != nil {
			return fmt.Errorf("failed to load layouts from %s: %w", source.Name, err)
		}
		// Load components
		if err := m.loadDirectory(source, ComponentsDir); err != nil {
			return fmt.Errorf("failed to load components from %s: %w", source.Name, err)
		}
		// Load partials
		if err := m.loadDirectory(source, PartialsDir); err != nil {
			return fmt.Errorf("failed to load partials from %s: %w", source.Name, err)
		}
	}

	return nil
}

// loadDirectory walks an entire directory tree looking for templates
func (m *Manager) loadDirectory(source TemplateSource, rootDir string) error {
	return fs.WalkDir(source.FS, rootDir, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return nil // Skip if directory doesn't exist
			}
			return fmt.Errorf("walk error for %s: %w", filePath, err)
		}

		if d.IsDir() {
			return nil
		}

		format := formatFromFile(filePath)
		if format == "" {
			return nil // Skip non-template files
		}

		// Read template content
		content, err := fs.ReadFile(source.FS, filePath)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", filePath, err)
		}

		// Parse into appropriate base template
		// Use the relative path from rootDir as the template name
		name := m.templateName(rootDir, filePath)
		base := m.baseTemplates[format]
		if _, err := base.New(name).Parse(string(content)); err != nil {
			return fmt.Errorf("failed to parse %s: %w", filePath, err)
		}

		return nil
	})
}

// templateName generates the template name from the root directory and file path
func (m *Manager) templateName(rootDir, filePath string) string {
	// Remove root directory prefix and extension
	name := strings.TrimPrefix(filePath, rootDir)
	name = strings.TrimPrefix(name, "/") // Remove leading slash if present
	name = strings.TrimSuffix(name, path.Ext(name))

	// Add prefix based on root directory
	switch rootDir {
	case LayoutsDir:
		return "layout:" + name
	case ComponentsDir:
		return "component:" + name
	case PartialsDir:
		return "partial:" + name
	default:
		return name
	}
}

// RenderedEmail represents a rendered email
type RenderedEmail struct {
	Subject string
	Text    string
	HTML    string
}

// RenderEmail renders an email template with optional layout
func (m *Manager) RenderEmail(name string, data interface{}, layout string) (*RenderedEmail, error) {
	if layout == "" {
		layout = m.defaultLayout
	}

	email := &RenderedEmail{}

	// Try text version
	if tmpl, err := m.getEmailTemplate(name, layout, FormatText); err == nil {
		text, err := m.executeTemplate(tmpl, "layout", data)
		if err != nil {
			return nil, fmt.Errorf("failed to render text template: %w", err)
		}
		email.Text = text
	}

	// Try HTML version
	if tmpl, err := m.getEmailTemplate(name, layout, FormatHTML); err == nil {
		html, err := m.executeTemplate(tmpl, "layout", data)
		if err != nil {
			return nil, fmt.Errorf("failed to render HTML template: %w", err)
		}

		if m.processor != nil {
			html, err = m.processor.Process(html)
			if err != nil {
				return nil, fmt.Errorf("failed to process HTML: %w", err)
			}
		}
		email.HTML = html
	}

	if email.Text == "" && email.HTML == "" {
		return nil, fmt.Errorf("no templates found for email %q", name)
	}

	return email, nil
}

// getEmailTemplate gets or creates an email template
func (m *Manager) getEmailTemplate(name, layout string, format TemplateFormat) (*template.Template, error) {
	cacheKey := fmt.Sprintf("%s:%s:%s", format, name, layout)

	m.mu.RLock()
	if tmpl, ok := m.emailCache[cacheKey]; ok {
		m.mu.RUnlock()
		return tmpl, nil
	}
	m.mu.RUnlock()

	// Need write lock to create template
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check cache again
	if tmpl, ok := m.emailCache[cacheKey]; ok {
		return tmpl, nil
	}

	// Clone base template
	base := m.baseTemplates[format]
	tmpl, err := base.Clone()
	if err != nil {
		return nil, err
	}

	// Find email template in sources (last one wins)
	filename := path.Join(EmailsDir, name+format.Extension())
	found := false

	for i := len(m.sources) - 1; i >= 0; i-- {
		source := m.sources[i]
		if content, err := fs.ReadFile(source.FS, filename); err == nil {
			if _, err := tmpl.New(name).Parse(string(content)); err != nil {
				return nil, err
			}
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("template %s not found", filename)
	}

	// Cache and return
	m.emailCache[cacheKey] = tmpl
	return tmpl, nil
}

// Extension returns the file extension for a template format
func (f TemplateFormat) Extension() string {
	switch f {
	case FormatHTML:
		return ".html"
	case FormatText:
		return ".txt"
	default:
		return ""
	}
}

// executeTemplate executes a template with the given name and data
func (m *Manager) executeTemplate(t *template.Template, name string, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, name, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ClearCache clears the email template cache
func (m *Manager) ClearCache() {
	m.mu.Lock()
	m.emailCache = make(map[string]*template.Template)
	m.mu.Unlock()
}

// AddFunc adds a function to the templates manager
func (m *Manager) AddFunc(name string, fn interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.funcMap == nil {
		m.funcMap = make(template.FuncMap)
	}

	m.funcMap[name] = fn
	return nil
}

// AddFuncs adds multiple functions to the templates manager
func (m *Manager) AddFuncs(funcs template.FuncMap) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.funcMap == nil {
		m.funcMap = make(template.FuncMap)
	}

	for name, fn := range funcs {
		m.funcMap[name] = fn
	}

	return nil
}

// themeFuncs returns the theme functions
func (m *Manager) themeFuncs() template.FuncMap {
	return template.FuncMap{
		"theme": func(path string) any {
			return GetThemeValue(m.theme, path)
		},
	}
}

// AddSource adds a new template source to the manager
func (m *Manager) AddSource(source TemplateSource) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Add to sources (later sources override earlier ones)
	m.sources = append(m.sources, source)

	// Clear cache since we have new sources
	m.emailCache = make(map[string]*template.Template)

	// Reload base templates to incorporate new source
	return m.loadBaseTemplates()
}
