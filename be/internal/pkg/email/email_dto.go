package email

// EmailConfig represents SMTP email configuration.
type EmailConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	From     string
}

// EmailRequest represents an email to be sent.
type EmailRequest struct {
	To      []string
	Subject string
	Body    string
	CC      []string
	BCC     []string
}
