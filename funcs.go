package mailpen

import (
	"fmt"
	"html/template"
)

// MergeFuncMaps merges the provided function maps into a single function map.
func MergeFuncMaps(maps ...template.FuncMap) template.FuncMap {
	result := make(template.FuncMap)
	for _, m := range maps {
		if m == nil {
			continue
		}
		for key, value := range m {
			result[key] = value
		}
	}
	return result
}

// cachedFuncMap holds the cached function map for the templates package.
var cachedFuncMap template.FuncMap

// DefaultFuncMap returns the complete function map for the templates package.
func DefaultFuncMap() template.FuncMap {
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
		"map_new": newMap, // Create a new map from key-value pairs
		"dict":    newMap, // Alias for map_new
		"add":     intAdd,
		"num_add": intAdd,
		"num_mod": mod,
		"sub":     intSub,
		"last":    indexLast,
	}
}

// intAdd adds two integers
func intAdd(a, b int) int {
	return a + b
}

// mod returns the remainder of a divided by b
func mod(a, b int) int {
	if b == 0 {
		return 0
	}
	return a % b
}

// intSub subtracts two integers
func intSub(a, b int) int {
	return a - b
}

// indexLast returns true if the index is the last element in the array
func indexLast(index int, arr []any) bool {
	return index == len(arr)-1
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
