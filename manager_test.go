package mailpen_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/patrickward/mailpen"
)

func TestManager_RenderEmail(t *testing.T) {
	tests := []struct {
		name        string
		sources     []mailpen.TemplateSource
		template    string
		layout      string
		data        map[string]any
		verify      func(*testing.T, *mailpen.RenderedEmail)
		wantErr     bool
		errContains string
	}{
		{
			name: "simple email uses the default layout",
			sources: []mailpen.TemplateSource{
				{
					Name: "default",
					FS:   testFS(t, "default"),
				},
			},
			template: "simple",
			data: map[string]any{
				"Name": "John Doe",
			},
			verify: func(t *testing.T, email *mailpen.RenderedEmail) {
				assert.Contains(t, email.HTML, "default-base-layout")
				assert.Contains(t, email.HTML, "Default HTML email without layout")
				assert.Contains(t, email.Text, "Default text email without layout")
			},
		},
		{
			name: "complete email with base override layout",
			sources: []mailpen.TemplateSource{
				{
					Name: "base",
					FS:   testFS(t, "base"),
				},
			},
			template: "welcome",
			data: map[string]any{
				"CompanyName":  "ACME Corp",
				"Name":         "John Doe",
				"LogoURL":      "https://example.com/logo.png",
				"SupportEmail": "support@example.com",
				"CurrentYear":  2024,
			},
			verify: func(t *testing.T, email *mailpen.RenderedEmail) {
				assert.Contains(t, email.HTML, "base-override-layout")
				assert.Contains(t, email.HTML, "<title>Welcome to ACME Corp</title>")
				assert.Contains(t, email.HTML, "Welcome, John Doe!")
				assert.Contains(t, email.HTML, "ACME Corp")
				assert.Contains(t, email.HTML, "support@example.com")
				assert.Contains(t, email.HTML, "2024")

				assert.Contains(t, email.Text, "Welcome, John Doe!")
				assert.Contains(t, email.Text, "ACME Corp")
				assert.Contains(t, email.Text, "support@example.com")
				assert.Contains(t, email.Text, "2024")
			},
		},
		{
			name: "override templates with second source",
			sources: []mailpen.TemplateSource{
				{
					Name: "base",
					FS:   testFS(t, "base"),
				},
				{
					Name: "override",
					FS:   testFS(t, "override"),
				},
			},
			template: "welcome",
			data: map[string]any{
				"CompanyName": "Override Corp",
				"Name":        "Jane Smith",
			},
			verify: func(t *testing.T, email *mailpen.RenderedEmail) {
				assert.Contains(t, email.HTML, "OVERRIDE Override Corp")
				assert.Contains(t, email.HTML, "Jane Smith")
				assert.Contains(t, email.Text, "OVERRIDE Override Corp")
			},
		},
		{
			name: "email with marketing layout",
			sources: []mailpen.TemplateSource{
				{
					Name: "base",
					FS:   testFS(t, "base"),
				},
			},
			template: "welcome",
			layout:   "marketing",
			data: map[string]any{
				"CompanyName": "ACME Corp",
				"Name":        "John Doe",
			},
			verify: func(t *testing.T, email *mailpen.RenderedEmail) {
				assert.Contains(t, email.HTML, `class="marketing-override-layout"`)
				assert.Contains(t, email.Text, "***") // Marketing text layout marker
			},
		},
		{
			name: "template not found",
			sources: []mailpen.TemplateSource{
				{
					Name: "base",
					FS:   testFS(t, "base"),
				},
			},
			template:    "nonexistent",
			wantErr:     true,
			errContains: "no templates found for email \"nonexistent\"",
		},
		{
			name: "invalid template syntax",
			sources: []mailpen.TemplateSource{
				{
					Name: "invalid",
					FS:   testFS(t, "invalid"),
				},
			},
			template:    "invalid",
			wantErr:     true,
			errContains: `no such template "@doesnotexist"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := mailpen.NewManager(&mailpen.ManagerConfig{
				Sources: tt.sources,
			})
			require.NoError(t, err)

			email, err := manager.RenderEmail(tt.template, tt.data, tt.layout)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			require.NoError(t, err)
			require.NotNil(t, email)

			if tt.verify != nil {
				tt.verify(t, email)
			}
		})
	}
}

func TestManager_AddSource(t *testing.T) {
	// Start with base templates
	manager, err := mailpen.NewManager(&mailpen.ManagerConfig{
		Sources: []mailpen.TemplateSource{
			{
				Name: "base",
				FS:   testFS(t, "base"),
			},
		},
	})
	require.NoError(t, err)

	// Verify base template rendering
	email, err := manager.RenderEmail("welcome", map[string]any{
		"CompanyName": "Base Corp",
		"Name":        "John Doe",
	}, "")
	require.NoError(t, err)
	assert.Contains(t, email.HTML, "Base Corp")

	// Add override source
	err = manager.AddSource(mailpen.TemplateSource{Name: "override", FS: testFS(t, "override")})
	require.NoError(t, err)

	// Verify override template is used
	email, err = manager.RenderEmail("welcome", map[string]any{
		"CompanyName": "Override Corp",
		"Name":        "Jane Smith",
	}, "")
	require.NoError(t, err)
	assert.Contains(t, email.HTML, "OVERRIDE Override Corp")
}

func TestManager_CacheClearing(t *testing.T) {
	manager, err := mailpen.NewManager(&mailpen.ManagerConfig{
		Sources: []mailpen.TemplateSource{
			{
				Name: "base",
				FS:   testFS(t, "base"),
			},
		},
	})
	require.NoError(t, err)

	// Render once to populate cache
	_, err = manager.RenderEmail("welcome", nil, "")
	require.NoError(t, err)

	// Clear cache
	manager.ClearCache()

	// Verify template still renders after cache clear
	email, err := manager.RenderEmail("welcome", map[string]any{
		"CompanyName": "ACME Corp",
		"Name":        "John Doe",
	}, "")
	require.NoError(t, err)
	assert.Contains(t, email.HTML, "Welcome, John Doe!")
}
