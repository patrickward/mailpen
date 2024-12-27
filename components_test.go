package mailpen_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/patrickward/mailpen"
)

// theme helper function for tests
func theme(path string) string {
	themeMap := mailpen.DefaultTheme()
	if val := mailpen.GetThemeValue(themeMap, path); val != nil {
		return val.(string)
	}
	return ""
}

func TestEmailComponents(t *testing.T) {
	// Create test configuration
	config := &mailpen.ManagerConfig{
		Sources: []mailpen.TemplateSource{
			{
				Name: "test",
				FS:   testFS(t, "base"),
			},
		},
	}

	manager, err := mailpen.NewManager(config)
	require.NoError(t, err)

	tests := []struct {
		name        string
		emailName   string
		data        map[string]interface{}
		wantHTML    []string // Strings that should be in HTML version
		wantText    []string // Strings that should be in text version
		notWantHTML []string // Strings that should not be in HTML version
	}{
		{
			name:      "email with headers",
			emailName: "headers-test",
			data: map[string]interface{}{
				"mainTitle":  "Welcome Message",
				"subTitle":   "Important Updates",
				"smallTitle": "Details Below",
			},
			wantHTML: []string{
				`<h1`,
				`color: #4DA647;`,
				`font-family: Arial, sans-serif;`,
				`font-size: 24px;`,
				`Welcome Message`,
				`</h1>`,
				`<h2`,
				`font-size: 18px;`,
				`Important Updates`,
				`</h2>`,
				`<h3`,
				`font-size: 16px;`,
				`Details Below`,
				`</h3>`,
			},
			wantText: []string{
				"Welcome Message",
				"Important Updates",
				"Details Below",
			},
		},
		{
			name:      "email with alert",
			emailName: "alert-test",
			data: map[string]interface{}{
				"alertTitle":   "Warning",
				"alertMessage": "Your account needs attention",
				"alertButton":  "Fix Now",
				"alertURL":     "https://example.com/fix",
			},
			wantHTML: []string{
				`border-left: 4px solid #ffa500;`,
				`color: #ffa500;`,
				`Warning`,
				`Your account needs attention`,
				`href="https://example.com/fix"`,
				`Fix Now`,
			},
			wantText: []string{
				"Warning",
				"Your account needs attention",
				"Fix Now",
				"https://example.com/fix",
			},
		},
		{
			name:      "email with logo",
			emailName: "logo-test",
			data: map[string]interface{}{
				"logoSrc": "/img/logo.png",
				"logoAlt": "Company Logo",
				"logoURL": "https://example.com",
			},
			wantHTML: []string{
				`src="/img/logo.png"`,
				`alt="Company Logo"`,
				`href="https://example.com"`,
			},
			wantText: []string{
				"Company Logo", // Alt text should be present in text version
			},
		},

		{
			name:      "email with data table",
			emailName: "table-test",
			data: map[string]interface{}{
				"tableData": mailpen.TableData{
					Headers: []mailpen.TableHeader{
						{Text: "Name", Width: "30%"},
						{Text: "Role", Width: "20%"},
						{Text: "Department", Width: "50%"},
					},
					Rows: []mailpen.TableRow{
						{
							Cells: []mailpen.TableCell{
								{Text: "John Doe", Width: "30%"},
								{Text: "Engineer", Width: "20%"},
								{Text: "Development", Width: "50%"},
							},
						},
						{
							Cells: []mailpen.TableCell{
								{Text: "Jane Smith", Width: "30%"},
								{Text: "Manager", Width: "20%"},
								{Text: "Operations", Width: "50%"},
							},
						},
					},
				},
			},
			wantHTML: []string{
				`<th`,
				`Name`,
				`Role`,
				`Department`,
				`John Doe`,
				`Engineer`,
				`Development`,
				`Jane Smith`,
				`Manager`,
				`Operations`,
				`font-family: Arial, sans-serif;`,
				`background-color: #4DA647;`,
				`width: 30%`,
				`width: 50%`,
			},
			wantText: []string{
				"Name |", "Role |", "Department",
				"John Doe |", "Engineer |", "Development",
				"Jane Smith |", "Manager |", "Operations",
			},
		},
		{
			name:      "email with card grid",
			emailName: "card-grid-test",
			data: map[string]interface{}{
				"cardData": mailpen.CardGridData{
					Cards: []mailpen.Card{
						{
							ImageURL:    "/images/product1.jpg",
							ImageAlt:    "Product One",
							Title:       "First Product",
							Description: "Description of first product",
							LinkURL:     "https://example.com/product1",
							LinkText:    "Learn More",
						},
						{
							ImageURL:    "/images/product2.jpg",
							ImageAlt:    "Product Two",
							Title:       "Second Product",
							Description: "Description of second product",
							LinkURL:     "https://example.com/product2",
							LinkText:    "Learn More",
						},
					},
				},
			},
			wantHTML: []string{
				`src="/images/product1.jpg"`,
				`alt="Product One"`,
				`First Product`,
				`Description of first product`,
				`href="https://example.com/product1"`,
				`Learn More`,
				`src="/images/product2.jpg"`,
				`Second Product`,
			},
			wantText: []string{
				"First Product",
				"Description of first product",
				"Learn More: https://example.com/product1",
				"Second Product",
				"Description of second product",
				"Learn More: https://example.com/product2",
			},
		},
		{
			name:      "email with buttons",
			emailName: "button-test",
			data: map[string]interface{}{
				"primaryButton": map[string]interface{}{
					"URL":  "https://example.com/primary",
					"Text": "Get Started",
				},
				"successButton": map[string]interface{}{
					"URL":   "https://example.com/success",
					"Text":  "Confirm",
					"Style": "success",
				},
				"dangerButton": map[string]interface{}{
					"URL":   "https://example.com/danger",
					"Text":  "Delete",
					"Style": "danger",
				},
			},
			wantHTML: []string{
				// Button links and text
				`href="https://example.com/primary"`,
				`Get Started`,
				`href="https://example.com/success"`,
				`Confirm`,
				`href="https://example.com/danger"`,
				`Delete`,
				// Style checks
				`background-color: ` + theme("colors.primary"),
				`background-color: ` + theme("colors.success"),
				`background-color: ` + theme("colors.danger"),
				`color: ` + theme("colors.background.primary"),
				`font-family: ` + theme("typography.font.family"),
				`border-radius: ` + theme("borders.radius.md"),
				`padding: ` + theme("components.button.padding.y") + ` ` + theme("components.button.padding.x"),
				`text-transform: ` + theme("components.button.textTransform"),
				`letter-spacing: ` + theme("typography.font.letterSpacing"),
			},
			wantText: []string{
				"Get Started: https://example.com/primary",
				"Confirm: https://example.com/success",
				"Delete: https://example.com/danger",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := manager.RenderEmail(tt.emailName, tt.data, "")
			require.NoError(t, err)

			// Check HTML content
			for _, want := range tt.wantHTML {
				assert.Contains(t, result.HTML, want)
			}
			for _, notWant := range tt.notWantHTML {
				assert.NotContains(t, result.HTML, notWant)
			}

			// Check text content
			for _, want := range tt.wantText {
				assert.Contains(t, result.Text, want)
			}
		})
	}
}
