package mailpen

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
)

// Message represents the content and recipients of an email message
type Message struct {
	From        string         // Sender email address
	To          []string       // List of recipient email addresses
	Cc          []string       // List of CC email addresses
	Bcc         []string       // List of BCC email addresses
	ReplyTo     string         // Reply-to email address
	Subject     string         // Email subject
	Data        map[string]any // Data to be passed to the templates
	Layout      string         // Layout name to process
	Template    string         // Template name to process
	TextBody    string         // Text body of the email
	HTMLBody    string         // HTML body of the email
	Attachments []Attachment   // List of attachments
}

// Attachment represents an email attachment
type Attachment struct {
	Filename    string
	Data        io.Reader
	ContentType ContentType
}

// Builder provides a fluent interface for constructing emails
type Builder struct {
	msg *Message
	err error
}

// NewMessage creates an email builder
func NewMessage() *Builder {
	return &Builder{
		msg: &Message{
			To:  make([]string, 0),
			Cc:  make([]string, 0),
			Bcc: make([]string, 0),
		},
	}
}

func (b *Builder) From(address string) *Builder {
	if b.err != nil {
		return b
	}
	b.msg.From = address
	return b
}

func (b *Builder) To(addresses ...string) *Builder {
	if b.err != nil {
		return b
	}
	b.msg.To = append(b.msg.To, addresses...)
	return b
}

func (b *Builder) Cc(addresses ...string) *Builder {
	if b.err != nil {
		return b
	}
	b.msg.Cc = append(b.msg.Cc, addresses...)
	return b
}

func (b *Builder) Bcc(addresses ...string) *Builder {
	if b.err != nil {
		return b
	}
	b.msg.Bcc = append(b.msg.Bcc, addresses...)
	return b
}

func (b *Builder) ReplyTo(address string) *Builder {
	if b.err != nil {
		return b
	}
	b.msg.ReplyTo = address
	return b
}

func (b *Builder) Subject(subject string) *Builder {
	if b.err != nil {
		return b
	}
	b.msg.Subject = subject
	return b
}

func (b *Builder) WithData(data map[string]any) *Builder {
	if b.err != nil {
		return b
	}
	b.msg.Data = data
	return b
}

func (b *Builder) Template(name string) *Builder {
	if b.err != nil {
		return b
	}
	b.msg.Template = name
	return b
}

func (b *Builder) Layout(name string) *Builder {
	if b.err != nil {
		return b
	}
	b.msg.Layout = name
	return b
}

// Attach adds an attachment to the email. The data is read from the provided reader and the content type is inferred from the filename.
func (b *Builder) Attach(filename string, data io.Reader) *Builder {
	if b.err != nil {
		return b
	}
	b.msg.Attachments = append(b.msg.Attachments, Attachment{
		Filename: filename,
		Data:     data,
	})
	return b
}

// AttachWithContentType adds an attachment to the email with a specific content type. The data is read from the provided reader.
func (b *Builder) AttachWithContentType(filename string, data io.Reader, contentType ContentType) *Builder {
	if b.err != nil {
		return b
	}
	b.msg.Attachments = append(b.msg.Attachments, Attachment{
		Filename:    filename,
		Data:        data,
		ContentType: contentType,
	})
	return b
}

// OpenFileAttachment is a helper that returns a file reader and a cleanup function
// for an attachment file. The filename is extracted from the filepath.
// It returns the filename, a reader for the file, a cleanup function, and an error if the file cannot be opened.
// It is the caller's responsibility to close the file reader after sending the email using the cleanup function.
//
// Example:
//
// filename, reader, cleanup, err := OpenFileAttachment("path/to/file.txt")
//
//	if err != nil {
//	    return err
//	}
//
// defer cleanup()
//
// msg, err := NewMessage().To("foo@example.com").Template("templates.tmpl").Attach(filename, reader).Build()
func OpenFileAttachment(filepath string) (string, io.Reader, func() error, error) {
	// If this is a directory, return an error indicating that directories are not supported
	info, err := os.Stat(filepath)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to get attachment info: %w", err)
	}

	if info.IsDir() {
		return "", nil, nil, errors.New("%s is a directory and directories are not supported as attachments")
	}

	f, err := os.Open(filepath)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to open attachment: %w", err)
	}

	filename := path.Base(filepath)
	cleanup := func() error {
		return f.Close()
	}

	return filename, f, cleanup, nil
}

func (b *Builder) Build() (*Message, error) {
	if b.err != nil {
		return nil, b.err
	}
	if len(b.msg.To) == 0 {
		return nil, errors.New("email must have at least one recipient")
	}
	return b.msg, nil
}

func (b *Builder) Must() *Message {
	msg, err := b.Build()
	if err != nil {
		panic(err)
	}
	return msg
}
