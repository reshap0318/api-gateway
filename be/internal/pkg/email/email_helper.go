package email

import (
	"encoding/base64"
	"fmt"
	"mime"
	"net/mail"
	"strings"
)

// ValidateEmail validates email format.
func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// BuildEmailMessage builds a complete email message with proper headers and encoding.
func BuildEmailMessage(config *EmailConfig, req EmailRequest) (string, error) {
	// Validate recipient emails (prevent CRLF injection)
	for _, email := range req.To {
		if !ValidateEmail(email) {
			return "", fmt.Errorf("invalid email: %s", email)
		}
	}

	var builder strings.Builder

	// From header (with optional name)
	fromName, fromEmail := parseFromAddress(config.From)
	if fromName != "" {
		builder.WriteString(fmt.Sprintf("From: %s <%s>\r\n", mime.QEncoding.Encode("utf-8", fromName), fromEmail))
	} else {
		builder.WriteString(fmt.Sprintf("From: %s\r\n", fromEmail))
	}

	// To header
	builder.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(req.To, ", ")))

	// CC header (if any)
	if len(req.CC) > 0 {
		builder.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(req.CC, ", ")))
	}

	// Subject (with encoding if needed)
	builder.WriteString(fmt.Sprintf("Subject: %s\r\n", encodeSubject(req.Subject)))

	// MIME headers
	builder.WriteString("MIME-Version: 1.0\r\n")
	builder.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	builder.WriteString("Content-Transfer-Encoding: base64\r\n")
	builder.WriteString("\r\n")

	// Body (base64 encoded)
	body := base64.StdEncoding.EncodeToString([]byte(req.Body))
	builder.WriteString(body)

	return builder.String(), nil
}

// parseFromAddress parses "Name <email>" format.
func parseFromAddress(from string) (name, email string) {
	addr, err := mail.ParseAddress(from)
	if err != nil {
		return "", from
	}
	return addr.Name, addr.Address
}

// encodeSubject encodes subject if contains non-ASCII characters using RFC 2047.
func encodeSubject(subject string) string {
	// Check if contains non-ASCII characters
	for i := 0; i < len(subject); i++ {
		if subject[i] > 127 {
			return mime.QEncoding.Encode("utf-8", subject)
		}
	}
	return subject
}
