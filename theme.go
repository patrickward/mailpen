package mailpen

import "strings"

// DefaultTheme returns a theme map that works with built-in templates
func DefaultTheme() map[string]any {
	return map[string]any{
		"colors": map[string]any{
			"primary":   "#4DA647", // From data-table header
			"secondary": "#30C3E6", // From header.html
			"success":   "#4caf50", // From button success
			"danger":    "#f44336", // From button danger
			"warning":   "#ffa500", // From button warning/default
			"text": map[string]any{
				"primary":   "#333333", // Dark text for main content
				"secondary": "#666666", // Used in cards and less prominent text
				"muted":     "#999999", // Used in footer text
			},
			"background": map[string]any{
				"primary":   "#ffffff",
				"secondary": "#f8f8f8", // Used in footer and quote backgrounds
			},
			"border": "#dddddd",
		},
		"typography": map[string]any{
			"font": map[string]any{
				"family": "Arial, sans-serif",
				"size": map[string]any{
					"xs":   "12px", // Footer text
					"sm":   "14px", // Secondary text
					"base": "16px", // Default body text
					"lg":   "18px", // Card titles
					"xl":   "24px", // Main headers
				},
				"lineHeight": map[string]any{
					"tight":   "18px",
					"normal":  "21px",
					"relaxed": "24px",
					"loose":   "30px",
				},
				"weight": map[string]any{
					"normal": "400",
					"medium": "500",
					"bold":   "700",
				},
				"letterSpacing": ".25px",
			},
		},
		"spacing": map[string]any{
			"0": "0",
			"1": "5px",
			"2": "10px",
			"3": "15px",
			"4": "20px", // Most common spacing
			"5": "30px",
			"6": "40px",
		},
		"borders": map[string]any{
			"width": "1px",
			"style": "solid",
			"radius": map[string]any{
				"sm": "3px",
				"md": "4px",
				"lg": "8px",
			},
		},
		"components": map[string]any{
			"button": map[string]any{
				"padding": map[string]any{
					"x": "24px",
					"y": "12px",
				},
				"textTransform": "uppercase",
			},
			"card": map[string]any{
				"padding": "20px",
				"shadow":  "none",
			},
			"table": map[string]any{
				"cell": map[string]any{
					"padding": "12px 15px",
				},
			},
			"notification": map[string]any{
				"padding":     "15px",
				"borderWidth": "4px",
			},
			"logo": map[string]any{
				"maxWidth": "200px",
				"padding":  "30px",
			},
		},
		"layout": map[string]any{
			"maxWidth": "600px",
			"gutter":   "20px",
		},
	}
}

// GetThemeValue safely traverses a theme map using dot notation
func GetThemeValue(theme map[string]any, path string) any {
	if path == "" {
		return nil
	}

	parts := strings.Split(path, ".")
	current := theme

	for i, part := range parts {
		if i == len(parts)-1 {
			return current[part]
		}

		next, ok := current[part].(map[string]any)
		if !ok {
			return nil
		}
		current = next
	}

	return nil
}
