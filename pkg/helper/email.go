package helper

import (
	"fmt"
	"net/smtp"
	"strings"
)

// SMTP struct defines the SMTP server configuration required to send emails.
type SMTP struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	FromName string `json:"from_name" yaml:"from_name"` // Sender's nickname
}

// SendEmailRequest defines the request parameters for sending an email.
type SendEmailRequest struct {
	To      []string // Recipients
	Cc      []string // Carbon copy
	Bcc     []string // Blind carbon copy
	Subject string   // Subject
	Body    string   // Content (HTML)
}

// SendEmail sends an email using the given configuration and request parameters.
func (s *SMTP) SendEmail(req SendEmailRequest) error {
	if len(req.To) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	// SMTP server address
	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)

	// Authentication information
	auth := smtp.PlainAuth("", s.User, s.Password, s.Host)

	// Build email headers
	headers := make(map[string]string)
	fromName := s.FromName
	if fromName == "" {
		fromName = s.User
	}
	headers["From"] = fmt.Sprintf("%s <%s>", fromName, s.User)
	headers["To"] = strings.Join(req.To, ", ")
	if len(req.Cc) > 0 {
		headers["Cc"] = strings.Join(req.Cc, ", ")
	}
	headers["Subject"] = req.Subject
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var msgBuilder strings.Builder
	for k, v := range headers {
		msgBuilder.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msgBuilder.WriteString("\r\n")
	msgBuilder.WriteString(req.Body)

	// Merge all recipients (To, Cc, Bcc)
	allRecipients := append(req.To, req.Cc...)
	allRecipients = append(allRecipients, req.Bcc...)

	// Send the email
	return smtp.SendMail(addr, auth, s.User, allRecipients, []byte(msgBuilder.String()))
}
