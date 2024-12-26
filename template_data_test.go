package mailpen_test

import (
	"reflect"
	"testing"

	"github.com/patrickward/mailpen"
)

func TestTemplateData_Merge(t *testing.T) {
	tests := []struct {
		name string
		base mailpen.TemplateData
		new  map[string]any
		want mailpen.TemplateData
	}{
		{
			name: "merge with empty base",
			base: mailpen.TemplateData{},
			new: map[string]any{
				"key": "value",
			},
			want: mailpen.TemplateData{
				"key": "value",
			},
		},
		{
			name: "merge with empty new data",
			base: mailpen.TemplateData{
				"existing": "value",
			},
			new: map[string]any{},
			want: mailpen.TemplateData{
				"existing": "value",
			},
		},
		{
			name: "merge overwrites existing values",
			base: mailpen.TemplateData{
				"key":  "old value",
				"keep": "kept value",
			},
			new: map[string]any{
				"key": "new value",
			},
			want: mailpen.TemplateData{
				"key":  "new value",
				"keep": "kept value",
			},
		},
		{
			name: "merge overwrites existing maps",
			base: mailpen.TemplateData{
				"links": map[string]any{
					"home":  "/home",
					"about": "/about",
				},
			},
			new: map[string]any{
				"links": map[string]any{
					"contact": "/contact",
				},
			},
			want: mailpen.TemplateData{
				"links": map[string]any{
					"contact": "/contact",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.base.Merge(tt.new)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Merge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTemplateData_MergeKeys(t *testing.T) {
	tests := []struct {
		name string
		base mailpen.TemplateData
		new  map[string]any
		want mailpen.TemplateData
	}{
		{
			name: "merge map into existing key",
			base: mailpen.TemplateData{
				"links": map[string]any{
					"home":  "/home",
					"about": "/about",
				},
			},
			new: map[string]any{
				"links": map[string]any{
					"contact": "/contact",
				},
			},
			want: mailpen.TemplateData{
				"links": map[string]any{
					"home":    "/home",
					"about":   "/about",
					"contact": "/contact",
				},
			},
		},
		{
			name: "add new top-level key",
			base: mailpen.TemplateData{
				"links": map[string]any{
					"home": "/home",
				},
			},
			new: map[string]any{
				"newKey": "value",
			},
			want: mailpen.TemplateData{
				"links": map[string]any{
					"home": "/home",
				},
				"newKey": "value",
			},
		},
		{
			name: "overwrite non-map value",
			base: mailpen.TemplateData{
				"key": "old value",
				"links": map[string]any{
					"home": "/home",
				},
			},
			new: map[string]any{
				"key": "new value",
				"links": map[string]any{
					"about": "/about",
				},
			},
			want: mailpen.TemplateData{
				"key": "new value",
				"links": map[string]any{
					"home":  "/home",
					"about": "/about",
				},
			},
		},
		{
			name: "handle non-map value overwriting map",
			base: mailpen.TemplateData{
				"key": map[string]any{
					"nested": "value",
				},
			},
			new: map[string]any{
				"key": "simple value",
			},
			want: mailpen.TemplateData{
				"key": "simple value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.base.MergeKeys(tt.new)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}
