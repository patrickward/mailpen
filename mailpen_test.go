package mailpen_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/patrickward/mailpen"
)

// mockProvider implements mailpen.Provider for testing
type mockProvider struct {
	sendCalls    int
	lastMessage  *mailpen.Message
	err          error
	capabilities mailpen.Capabilities
}

func (m *mockProvider) Send(ctx context.Context, msg *mailpen.Message) error {
	m.sendCalls++
	m.lastMessage = msg
	return m.err
}

func (m *mockProvider) Name() string {
	return "mock"
}

func (m *mockProvider) Validate(msg *mailpen.Message) error {
	if m.err != nil {
		return m.err
	}
	if len(msg.To) == 0 {
		return errors.New("at least one recipient is required")
	}
	return nil
}

func (m *mockProvider) Capabilities() mailpen.Capabilities {
	return m.capabilities
}

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		provider   mailpen.Provider
		config     *mailpen.Config
		opts       []mailpen.Option
		wantErr    bool
		errMessage string
	}{
		{
			name:       "nil provider",
			provider:   nil,
			config:     &mailpen.Config{},
			wantErr:    true,
			errMessage: "provider is required",
		},
		{
			name:       "nil config",
			provider:   &mockProvider{},
			config:     nil,
			wantErr:    true,
			errMessage: "config is required",
		},
		{
			name:     "valid creation",
			provider: &mockProvider{},
			config: &mailpen.Config{
				From: "test@example.com",
				Sources: []mailpen.TemplateSource{
					{
						Name: "base",
						FS:   testFS(t, "base"),
					},
				},
			},
		},
		{
			name:     "with override templates",
			provider: &mockProvider{},
			config: &mailpen.Config{
				From: "test@example.com",
				Sources: []mailpen.TemplateSource{
					{
						Name: "base",
						FS:   testFS(t, "base"),
					},
					{
						Name: "override",
						FS:   testFS(t, "override"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp, err := mailpen.New(tt.provider, tt.config, tt.opts...)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMessage != "" {
					assert.Contains(t, err.Error(), tt.errMessage)
				}
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, mp)
		})
	}
}

func TestMailpen_Send(t *testing.T) {
	tests := []struct {
		name       string
		config     *mailpen.Config
		message    *mailpen.Message
		mockSetup  func(*mockProvider)
		verify     func(*testing.T, *mockProvider)
		wantErr    bool
		errMessage string
	}{
		{
			name: "basic send without template",
			config: &mailpen.Config{
				From: "sender@example.com",
			},
			message: mailpen.NewMessage().
				To("recipient@example.com").
				Subject("Test Subject").
				WithData(map[string]any{"name": "John"}).
				Must(),
			verify: func(t *testing.T, m *mockProvider) {
				assert.Equal(t, 1, m.sendCalls)
				require.NotNil(t, m.lastMessage)
				assert.Equal(t, "recipient@example.com", m.lastMessage.To[0])
			},
		},
		{
			name: "send with template",
			config: &mailpen.Config{
				From: "sender@example.com",
				Sources: []mailpen.TemplateSource{
					{
						Name: "base",
						FS:   testFS(t, "base"),
					},
				},
			},
			message: mailpen.NewMessage().
				To("recipient@example.com").
				Template("welcome").
				WithData(map[string]any{
					"Name":        "John",
					"CompanyName": "ACME Corp",
				}).
				Must(),
			verify: func(t *testing.T, m *mockProvider) {
				assert.Equal(t, 1, m.sendCalls)
				require.NotNil(t, m.lastMessage)
				assert.Contains(t, m.lastMessage.HTMLBody, "Welcome, John!")
				assert.Contains(t, m.lastMessage.HTMLBody, "ACME Corp")
			},
		},
		{
			name: "send failure",
			config: &mailpen.Config{
				From: "sender@example.com",
			},
			message: mailpen.NewMessage().
				To("recipient@example.com").
				Subject("Test").
				Must(),
			mockSetup: func(m *mockProvider) {
				m.err = errors.New("send failed")
			},
			wantErr:    true,
			errMessage: "send failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProvider{}
			if tt.mockSetup != nil {
				tt.mockSetup(mock)
			}

			mp, err := mailpen.New(mock, tt.config)
			require.NoError(t, err)

			err = mp.Send(context.Background(), tt.message)
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
