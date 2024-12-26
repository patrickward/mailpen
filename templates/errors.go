package templates

import (
	"fmt"
)

// TemplateError provides context about templates errors
type TemplateError struct {
	TemplateName string
	OriginalErr  error
	Phase        string // "parse", "execute", "process"
}

// Error implements the error interface
func (e *TemplateError) Error() string {
	if e.TemplateName == "" {
		return fmt.Sprintf("template error during %s phase: %v", e.Phase, e.OriginalErr)
	}

	return fmt.Sprintf("template error in %s during %s phase: %v", e.TemplateName, e.Phase, e.OriginalErr)
}
