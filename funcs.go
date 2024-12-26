package mailpen

import (
	"fmt"
	"html/template"
)

// MergeFuncMaps merges the provided function maps into a single function map.
func MergeFuncMaps(maps ...template.FuncMap) template.FuncMap {
	result := make(template.FuncMap)
	for _, m := range maps {
		for key, value := range m {
			result[key] = value
		}
	}
	return result
}

// cachedFuncMap holds the cached function map for the templates package.
var cachedFuncMap template.FuncMap

// FuncMap returns the complete function map for the templates package.
func FuncMap() template.FuncMap {
	if cachedFuncMap != nil {
		return cachedFuncMap
	}

	// TODO: Add default function maps here
	cachedFuncMap = MergeFuncMaps(
		mapFuncs(),
	)

	return cachedFuncMap
}

// Helper functions for template functions
func mapFuncs() template.FuncMap {
	return template.FuncMap{
		"map_new": newMap,
	}
}

// newMap creates a new map from key-value pairs
//
// Example: {{ map.new "key" "value" "other" "value" }} -> map[key:value other:value]
func newMap(pairs ...any) (map[string]any, error) {
	if len(pairs)%2 != 0 {
		return nil, fmt.Errorf("map.new requires pairs of arguments")
	}

	result := make(map[string]any)
	for i := 0; i < len(pairs); i += 2 {
		key, ok := pairs[i].(string)
		if !ok {
			return nil, fmt.Errorf("map key must be string, got %T", pairs[i])
		}
		result[key] = pairs[i+1]
	}
	return result, nil
}
