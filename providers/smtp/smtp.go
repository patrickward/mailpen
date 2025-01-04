package smtp

import (
	"context"
	"fmt"
	"time"

	gomail "github.com/wneessen/go-mail"

	"github.com/patrickward/mailpen"
)

// Client defines the interface for an SMTP client
type Client interface {
	DialAndSend(messages ...*gomail.Msg) error
}

// Config holds SMTP-specific configuration
type Config struct {
	Host      string
	Port      int
	Username  string
	Password  string
	AuthType  string // Type of SMTP authentication
	TLSPolicy int    // TLS policy for the SMTP connection

	// Retry configuration
	RetryCount int
	RetryDelay time.Duration
}

type Provider struct {
	client Client
	config *Config
}

type Option func(p *Provider)

// WithClient allows injection of a custom SMTP client
func WithClient(client Client) Option {
	return func(p *Provider) {
		p.client = client
	}
}

// New creates a new SMTP provider
func New(config *Config, opts ...Option) (*Provider, error) {
	if config == nil {
		return nil, fmt.Errorf("config is required")
	}

	if config.RetryCount == 0 {
		config.RetryCount = 1
	}

	authType := authTypeFromString(config.AuthType)
	tlsPolicy := tlsPolicyFromInt(config.TLSPolicy)

	client, err := gomail.NewClient(
		config.Host,
		gomail.WithTimeout(10*time.Second),
		gomail.WithSMTPAuth(authType),
		gomail.WithPort(config.Port),
		gomail.WithUsername(config.Username),
		gomail.WithPassword(config.Password),
		gomail.WithTLSPolicy(tlsPolicy),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create SMTP client: %w", err)
	}

	p := &Provider{
		client: client,
		config: config,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p, nil
}

// Send implements mailpen.Provider
func (p *Provider) Send(ctx context.Context, msg *mailpen.Message) error {
	email := gomail.NewMsg()
	email.Subject(msg.Subject)

	if err := p.setAddresses(email, msg); err != nil {
		return err
	}

	if err := p.setBodies(email, msg); err != nil {
		return err
	}

	if err := p.addAttachments(email, msg.Attachments); err != nil {
		return err
	}

	return p.sendWithRetry(email)
}

func (p *Provider) Name() string {
	return "smtp"
}

func (p *Provider) Validate(msg *mailpen.Message) error {
	if len(msg.To) == 0 {
		return fmt.Errorf("at least one recipient is required")
	}
	return nil
}

func (p *Provider) Capabilities() mailpen.Capabilities {
	return mailpen.Capabilities{
		MaxRecipients:      1000,
		MaxAttachmentSize:  25 * 1024 * 1024,
		SupportsTemplates:  true,
		SupportsHTMLOnly:   true,
		SupportsScheduling: false,
	}
}

// addAttachments adds attachments to the email
func (p *Provider) addAttachments(email *gomail.Msg, attachments []mailpen.Attachment) error {
	for _, att := range attachments {
		var opts []gomail.FileOption
		if att.ContentType != "" {
			opts = append(opts, gomail.WithFileContentType(toGoMailContentType(att.ContentType.String())))
		}

		if att.Data == nil {
			return fmt.Errorf("nil reader for attachment %s", att.Filename)
		}

		if err := email.AttachReader(att.Filename, att.Data, opts...); err != nil {
			return fmt.Errorf("failed to attach file %s: %w", att.Filename, err)
		}
	}
	return nil
}

// setAddresses sets the addresses on the email
func (p *Provider) setAddresses(email *gomail.Msg, msg *mailpen.Message) error {
	if err := email.From(msg.From); err != nil {
		return fmt.Errorf("failed to set from address: %w", err)
	}

	if err := email.To(msg.To...); err != nil {
		return fmt.Errorf("failed to set to addresses: %w", err)
	}

	if len(msg.Cc) > 0 {
		if err := email.Cc(msg.Cc...); err != nil {
			return fmt.Errorf("failed to set cc addresses: %w", err)
		}
	}

	if len(msg.Bcc) > 0 {
		if err := email.Bcc(msg.Bcc...); err != nil {
			return fmt.Errorf("failed to set bcc addresses: %w", err)
		}
	}

	if msg.ReplyTo != "" {
		if err := email.ReplyTo(msg.ReplyTo); err != nil {
			return fmt.Errorf("failed to set reply-to address: %w", err)
		}
	}

	return nil
}

// setBodies sets the text and HTML bodies on the email
func (p *Provider) setBodies(email *gomail.Msg, msg *mailpen.Message) error {
	if msg.TextBody != "" {
		email.SetBodyString(gomail.TypeTextPlain, msg.TextBody)
	}

	if msg.HTMLBody != "" {
		if msg.TextBody != "" {
			email.AddAlternativeString(gomail.TypeTextHTML, msg.HTMLBody)
		} else {
			email.SetBodyString(gomail.TypeTextHTML, msg.HTMLBody)
		}
	}

	return nil
}

// sendWithRetry sends the email with retries
func (p *Provider) sendWithRetry(email *gomail.Msg) error {
	var lastErr error
	for i := 0; i < p.config.RetryCount; i++ {
		if err := p.client.DialAndSend(email); err != nil {
			lastErr = err
			if i < p.config.RetryCount-1 {
				time.Sleep(p.config.RetryDelay)
				continue
			}
		} else {
			return nil
		}
	}
	return fmt.Errorf("failed to send email after %d attempts: %w", p.config.RetryCount, lastErr)
}

// authTypeFromString converts a string to a gomail.SMTPAuthType
func authTypeFromString(typ string) gomail.SMTPAuthType {
	switch typ {
	case "PLAIN":
		return gomail.SMTPAuthPlain
	case "LOGIN":
		return gomail.SMTPAuthLogin
	case "CRAM-MD5":
		return gomail.SMTPAuthCramMD5
	case "NOAUTH":
		return gomail.SMTPAuthNoAuth
	case "XOAUTH2":
		return gomail.SMTPAuthXOAUTH2
	case "CUSTOM":
		return gomail.SMTPAuthCustom
	case "SCRAM-SHA-1":
		return gomail.SMTPAuthSCRAMSHA1
	case "SCRAM-SHA-1-PLUS":
		return gomail.SMTPAuthSCRAMSHA1PLUS
	case "SCRAM-SHA-256":
		return gomail.SMTPAuthSCRAMSHA256
	case "SCRAM-SHA-256-PLUS":
		return gomail.SMTPAuthSCRAMSHA256PLUS
	default:
		return gomail.SMTPAuthLogin
	}
}

// tlsPolicyFromInt converts an integer to a gomail.TLSPolicy
func tlsPolicyFromInt(typ int) gomail.TLSPolicy {
	switch typ {
	case 0:
		return gomail.NoTLS
	case 1:
		return gomail.TLSOpportunistic
	case 2:
		return gomail.TLSMandatory
	default:
		return gomail.TLSOpportunistic
	}
}

// toGoMailContentType converts the ContentType to a string that can be used with the gomail package.
func toGoMailContentType(contentType string) gomail.ContentType {
	return gomail.ContentType(contentType)
}
