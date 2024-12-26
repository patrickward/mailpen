package smtp_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomail "github.com/wneessen/go-mail"

	"github.com/patrickward/mailpen"
	"github.com/patrickward/mailpen/providers/smtp"
)

// mockSMTPClient implements smtp.SMTPClient for testing
type mockSMTPClient struct {
	sendCalls int
	messages  []*gomail.Msg
	err       error
}

func (m *mockSMTPClient) DialAndSend(messages ...*gomail.Msg) error {
	m.sendCalls++
	if m.err != nil {
		return m.err
	}
	m.messages = append(m.messages, messages...)
	return nil
}

func TestProvider_Send(t *testing.T) {
	tests := []struct {
		name       string
		config     *smtp.Config
		message    *mailpen.Message
		mockSetup  func(*mockSMTPClient)
		verify     func(*testing.T, *mockSMTPClient)
		wantErr    bool
		errMessage string
	}{
		{
			name: "successful send",
			config: &smtp.Config{
				Host:     "smtp.example.com",
				Port:     587,
				Username: "test",
				Password: "pass",
			},
			message: &mailpen.Message{
				From:     "sender@example.com",
				To:       []string{"recipient@example.com"},
				Subject:  "Test Email",
				TextBody: "Hello World",
			},
			verify: func(t *testing.T, m *mockSMTPClient) {
				assert.Equal(t, 1, m.sendCalls)
				require.Len(t, m.messages, 1)
				msg := m.messages[0]

				from := msg.GetFrom()
				require.GreaterOrEqual(t, len(from), 1)
				assert.Equal(t, "sender@example.com", from[0].Address)

				to, err := msg.GetRecipients()
				require.NoError(t, err)
				require.GreaterOrEqual(t, len(to), 1)
				assert.Equal(t, []string{"recipient@example.com"}, to)
			},
		},
		{
			name: "retry on failure",
			config: &smtp.Config{
				Host:       "smtp.example.com",
				Port:       587,
				RetryCount: 3,
				RetryDelay: time.Millisecond,
			},
			message: &mailpen.Message{
				From:    "sender@example.com",
				To:      []string{"recipient@example.com"},
				Subject: "Test Email",
			},
			mockSetup: func(m *mockSMTPClient) {
				m.err = &gomail.SendError{}
			},
			verify: func(t *testing.T, m *mockSMTPClient) {
				assert.Equal(t, 3, m.sendCalls) // Verifies retry attempts
			},
			wantErr:    true,
			errMessage: "failed to send email after 3 attempts",
		},
		{
			name: "with attachments",
			config: &smtp.Config{
				Host: "smtp.example.com",
				Port: 587,
			},
			message: &mailpen.Message{
				From:    "sender@example.com",
				To:      []string{"recipient@example.com"},
				Subject: "Test Email",
				Attachments: []mailpen.Attachment{
					{
						Filename:    "test.txt",
						Data:        strings.NewReader("test content"),
						ContentType: mailpen.TypeTextPlain,
					},
				},
			},
			verify: func(t *testing.T, m *mockSMTPClient) {
				assert.Equal(t, 1, m.sendCalls)
				require.Len(t, m.messages, 1)
			},
		},
		{
			name: "with cc and bcc",
			config: &smtp.Config{
				Host: "smtp.example.com",
				Port: 587,
			},
			message: &mailpen.Message{
				From:    "sender@example.com",
				To:      []string{"recipient@example.com"},
				Cc:      []string{"cc@example.com"},
				Bcc:     []string{"bcc@example.com"},
				Subject: "Test Email",
			},
			verify: func(t *testing.T, m *mockSMTPClient) {
				require.Len(t, m.messages, 1)
				msg := m.messages[0]

				cc := msg.GetCc()
				require.GreaterOrEqual(t, len(cc), 1)
				assert.Equal(t, "cc@example.com", cc[0].Address)

				bcc := msg.GetBcc()
				require.GreaterOrEqual(t, len(bcc), 1)
				assert.Equal(t, "bcc@example.com", bcc[0].Address)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockSMTPClient{}
			if tt.mockSetup != nil {
				tt.mockSetup(mock)
			}

			provider, err := smtp.New(tt.config, smtp.WithClient(mock))
			require.NoError(t, err)

			err = provider.Send(context.Background(), tt.message)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMessage != "" {
					assert.Contains(t, err.Error(), tt.errMessage)
				}
				return
			}
			require.NoError(t, err)

			if tt.verify != nil {
				tt.verify(t, mock)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		config     *smtp.Config
		opts       []smtp.Option
		wantErr    bool
		errMessage string
	}{
		{
			name: "valid config",
			config: &smtp.Config{
				Host:     "smtp.example.com",
				Port:     587,
				Username: "test",
				Password: "pass",
			},
		},
		{
			name:       "nil config",
			config:     nil,
			wantErr:    true,
			errMessage: "config is required",
		},
		{
			name: "with custom client",
			config: &smtp.Config{
				Host: "smtp.example.com",
				Port: 587,
			},
			opts: []smtp.Option{
				smtp.WithClient(&mockSMTPClient{}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := smtp.New(tt.config, tt.opts...)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMessage != "" {
					assert.Contains(t, err.Error(), tt.errMessage)
				}
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, provider)
		})
	}
}
