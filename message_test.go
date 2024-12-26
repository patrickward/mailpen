package mailpen_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/patrickward/mailpen"
)

func TestMessageBuilder(t *testing.T) {
	tests := []struct {
		name      string
		build     func(*mailpen.Builder)
		wantErr   bool
		errString string
		validate  func(*testing.T, *mailpen.Message)
	}{
		{
			name: "basic message",
			build: func(b *mailpen.Builder) {
				b.To("user@example.com").
					WithData(map[string]any{"name": "John"})
			},
			validate: func(t *testing.T, msg *mailpen.Message) {
				assert.Equal(t, []string{"user@example.com"}, msg.To)
				assert.Equal(t, map[string]any{"name": "John"}, msg.Data)
			},
		},
		{
			name: "message with cc and bcc",
			build: func(b *mailpen.Builder) {
				b.To("user@example.com").
					Cc("cc@example.com").
					Bcc("bcc@example.com")
			},
			validate: func(t *testing.T, msg *mailpen.Message) {
				assert.Equal(t, []string{"user@example.com"}, msg.To)
				assert.Equal(t, []string{"cc@example.com"}, msg.Cc)
				assert.Equal(t, []string{"bcc@example.com"}, msg.Bcc)
			},
		},
		{
			name: "message with reply-to",
			build: func(b *mailpen.Builder) {
				b.To("user@example.com").
					ReplyTo("reply@example.com")
			},
			validate: func(t *testing.T, msg *mailpen.Message) {
				assert.Equal(t, []string{"user@example.com"}, msg.To)
				assert.Equal(t, "reply@example.com", msg.ReplyTo)
			},
		},
		{
			name:      "missing recipient",
			build:     func(b *mailpen.Builder) {},
			wantErr:   true,
			errString: "email must have at least one recipient",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := mailpen.NewMessage()
			tt.build(b)
			msg, err := b.Build()

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errString)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, msg)
			if tt.validate != nil {
				tt.validate(t, msg)
			}
		})
	}
}
