package email

import (
	"fmt"
	"net/smtp"

	"github.com/reshap0318/api-gateway/internal/helpers"
	"github.com/reshap0318/api-gateway/internal/pkg/email/templates"
)

// EmailClient represents an email client for sending emails via SMTP.
type EmailClient struct {
	config *EmailConfig
}

// NewEmailClient creates and initializes a new EmailClient.
func NewEmailClient() *EmailClient {
	return &EmailClient{
		config: &EmailConfig{
			Host:     helpers.GetEnv("SMTP_HOST", ""),
			Port:     helpers.GetEnv("SMTP_PORT", "587"),
			User:     helpers.GetEnv("SMTP_USER", ""),
			Password: helpers.GetEnv("SMTP_PASSWORD", ""),
			From:     helpers.GetEnv("SMTP_FROM", "noreply@example.com"),
		},
	}
}

// IsConfigured returns true if SMTP is configured.
func (c *EmailClient) IsConfigured() bool {
	return c.config.Host != "" && c.config.User != "" && c.config.Password != ""
}

// SendEmail sends an email using SMTP.
func (c *EmailClient) SendEmail(req EmailRequest) error {
	if !c.IsConfigured() {
		// Log warning but don't fail (for development)
		fmt.Println("[Email] SMTP not configured, skipping email send")
		return nil
	}

	if len(req.To) == 0 {
		return &helpers.FieldError{Field: "email", Message: "invalid email address"}
	}

	// Build email message
	message, err := BuildEmailMessage(c.config, req)
	if err != nil {
		return fmt.Errorf("failed to build email message: %w", err)
	}

	// Send to each recipient
	for _, to := range req.To {
		if err := c.send(to, message); err != nil {
			return fmt.Errorf("failed to send email to %s: %w", to, err)
		}
	}

	return nil
}

// SendResetPasswordEmail sends a reset password email.
func (c *EmailClient) SendResetPasswordEmail(to, token, resetURL string) error {
	appName := helpers.GetEnv("APP_NAME", "Our Application")
	
	body := templates.ResetPasswordEmail(token, resetURL, appName)

	req := EmailRequest{
		To:      []string{to},
		Subject: "Reset Password Request",
		Body:    body,
	}

	return c.SendEmail(req)
}

// send sends an email to a single recipient.
func (c *EmailClient) send(to, message string) error {
	addr := fmt.Sprintf("%s:%s", c.config.Host, c.config.Port)

	auth := smtp.PlainAuth(
		"",
		c.config.User,
		c.config.Password,
		c.config.Host,
	)

	return smtp.SendMail(
		addr,
		auth,
		c.config.From,
		[]string{to},
		[]byte(message),
	)
}
