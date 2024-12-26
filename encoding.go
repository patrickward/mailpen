package mailpen

// ContentType is a type wrapper for a string and represents the MIME type of the content being handled.
type ContentType string

const (
	// TypeAppOctetStream represents the MIME type for arbitrary binary data.
	TypeAppOctetStream ContentType = "application/octet-stream"

	// TypeMultipartAlternative represents the MIME type for a message body that can contain multiple alternative
	// formats.
	TypeMultipartAlternative ContentType = "multipart/alternative"

	// TypeMultipartMixed represents the MIME type for a multipart message containing different parts.
	TypeMultipartMixed ContentType = "multipart/mixed"

	// TypeMultipartRelated represents the MIME type for a multipart message where each part is a related file
	// or resource.
	TypeMultipartRelated ContentType = "multipart/related"

	// TypePGPSignature represents the MIME type for PGP signed messages.
	TypePGPSignature ContentType = "application/pgp-signature"

	// TypePGPEncrypted represents the MIME type for PGP encrypted messages.
	TypePGPEncrypted ContentType = "application/pgp-encrypted"

	// TypeTextHTML represents the MIME type for HTML text content.
	TypeTextHTML ContentType = "text/html"

	// TypeTextPlain represents the MIME type for plain text content.
	TypeTextPlain ContentType = "text/plain"
)

// String returns the string representation of the ContentType and implements the Stringer interface.
func (c ContentType) String() string {
	return string(c)
}
