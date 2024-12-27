# Mailpen Package Documentation - Experimental

⚠️ **EXPERIMENTAL**: This project is under active development and the API changes frequently. Not recommended for production use unless you're willing to vendor the code.

## Overview
Mailpen is a flexible email composition and sending system for Go applications. It provides a modular architecture for sending emails with support for templates, layouts, components, and multiple email providers.

## Core Concepts

### Message Structure
Messages in Mailpen consist of:
- Recipients (To, CC, BCC)
- Subject
- Content (HTML and/or Text)
- Optional attachments
- Optional template data

### Templates Organization
Templates are organized into four main directories:
- `layouts/` - Base email layouts
- `components/` - Reusable email components
- `partials/` - Shared template fragments
- `emails/` - Individual email templates

### Template Sources
Multiple template sources can be configured, allowing for template overrides and customization. Sources are processed in order, with later sources taking precedence.

## Basic Usage

```go
// Create a configuration
config := &mailpen.Config{
    From: "sender@example.com",
    CompanyName: "ACME Corp",
    Sources: []mailpen.TemplateSource{
        {
            Name: "base",
            FS: os.DirFS("templates/base"),
        },
    },
}

// Initialize a provider (e.g., SMTP)
provider := smtp.New(&smtp.Config{
    Host: "smtp.example.com",
    Port: 587,
})

// Create Mailpen instance
mp, err := mailpen.New(provider, config)
if err != nil {
    log.Fatal(err)
}

// Send an email using a template
msg := mailpen.NewMessage().
    To("recipient@example.com").
    Template("welcome").
    WithData(map[string]any{
        "Name": "John Doe",
    }).
    Must()

err = mp.Send(context.Background(), msg)
```

## Adding New Components

### 1. Create Component Template
Components should be added to the `components/` directory. Each component should have both HTML and text versions:

```
components/
  ├── alert/
  │   ├── alert.html
  │   └── alert.txt
  └── button/
      ├── button.html
      └── button.txt
```

### 2. Define Component Structure
Create a struct in `components.go` to define the component's data structure:

```go
type AlertData struct {
    Title    string
    Message  string
    Type     string
    ButtonURL string
}
```

### 3. Implement Component Templates
HTML template example (`components/alert/alert.html`):
```html
{{define "component:alert"}}
<div class="alert alert-{{.Type}}" style="
    border-left: 4px solid {{theme "colors.warning"}};
    padding: {{theme "components.notification.padding"}};
">
    <h4>{{.Title}}</h4>
    <p>{{.Message}}</p>
    {{if .ButtonURL}}
    <a href="{{.ButtonURL}}" class="btn">Fix Now</a>
    {{end}}
</div>
{{end}}
```

Text template example (`components/alert/alert.txt`):
```
{{define "component:alert"}}
! {{.Title}} !
{{.Message}}
{{if .ButtonURL}}
Fix Now: {{.ButtonURL}}
{{end}}
{{end}}
```

## Adding New Layouts

### 1. Create Layout Files
Layouts should be added to the `layouts/` directory with both HTML and text versions:

```
layouts/
  ├── base/
  │   ├── base.html
  │   └── base.txt
  └── marketing/
      ├── marketing.html
      └── marketing.txt
```

### 2. Implement Layout Structure
HTML layout example (`layouts/marketing/marketing.html`):
```html
{{define "layout:marketing"}}
<!DOCTYPE html>
<html>
<head>
    <title>{{.Subject}}</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: {{theme "typography.font.family"}};
            line-height: {{theme "typography.font.lineHeight.normal"}};
            color: {{theme "colors.text.primary"}};
        }
    </style>
</head>
<body>
    <div class="container" style="max-width: {{theme "layout.maxWidth"}}">
        {{template "content" .}}
        {{template "component:footer" .FooterData}}
    </div>
</body>
</html>
{{end}}
```

Text layout example (`layouts/marketing/marketing.txt`):
```
{{define "layout:marketing"}}
{{template "content" .}}

---
{{template "component:footer" .FooterData}}
{{end}}
```

## Theming

Mailpen includes a theming system that can be customized through the configuration. Theme values can be accessed in templates using the `theme` function:

```html
<div style="color: {{theme "colors.primary"}}">
    Themed content
</div>
```

Default theme values include:
- Colors (primary, secondary, success, danger, warning)
- Typography (font families, sizes, weights)
- Spacing
- Border styles
- Component-specific styles

## Best Practices

1. **Template Organization**
    - Keep templates modular and reusable
    - Use components for repeated elements
    - Maintain consistent naming conventions

2. **Theme Usage**
    - Use theme values for consistent styling
    - Avoid hardcoding colors or dimensions
    - Define new theme values for custom components

3. **Testing**
    - Test both HTML and text renderings
    - Verify template data handling
    - Check responsive layouts

4. **Error Handling**
    - Always check template parsing errors
    - Validate required fields
    - Handle attachment errors properly

## Common Issues

1. **Template Not Found**
    - Verify template path matches the name used in code
    - Check template source order for overrides
    - Ensure both HTML and text versions exist

2. **Style Inconsistencies**
    - Use theme values consistently
    - Test email rendering in multiple clients
    - Follow email HTML best practices

3. **Missing Data**
    - Validate required template data
    - Provide default values where appropriate
    - Use template functions for data formatting

## Advanced Features

### Custom Template Functions
Add custom template functions using `AddFunc` or `AddFuncs`:

```go
manager.AddFunc("formatDate", func(t time.Time) string {
    return t.Format("2006-01-02")
})
```

### Multiple Template Sources
Configure multiple sources for template overrides:

```go
config := &mailpen.Config{
    Sources: []mailpen.TemplateSource{
        {Name: "base", FS: baseFS},
        {Name: "custom", FS: customFS},
    },
}
```

### HTML Processing
Implement custom HTML processing:

```go
type CustomProcessor struct{}

func (p *CustomProcessor) Process(html string) (string, error) {
    // Process HTML
    return html, nil
}

config.HTMLProcessor = &CustomProcessor{}
```
