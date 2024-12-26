package mailpen

import (
	"context"
)

// Provider defines the interface for email providers
type Provider interface {
	// Send sends an email message
	Send(ctx context.Context, msg *Message) error

	// Name returns the provider name
	Name() string

	// Validate validates a message before sending
	Validate(msg *Message) error

	// Capabilities returns the provider capabilities
	Capabilities() Capabilities
}

// Capabilities defines what features a provider supports
type Capabilities struct {
	MaxRecipients      int
	MaxAttachmentSize  int64
	SupportsTemplates  bool
	SupportsHTMLOnly   bool
	SupportsScheduling bool
}
