package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"own-paynet/config"
	"path/filepath"
)

type EmailService struct {
	Config *config.Config
}

type EmailData struct {
	To       string
	Subject  string
	Template string
	Data     map[string]interface{}
}

// NewEmailService creates a new email service
func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{Config: cfg}
}

// SendEmail sends an email using the specified template and data
func (s *EmailService) SendEmail(data EmailData) error {
	// Parse template
	tmpl, err := template.ParseFiles(filepath.Join("utils/email/templates", data.Template))
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	// Execute template with data
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data.Data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	// Set up authentication information
	auth := smtp.PlainAuth(
		"",
		s.Config.SMTPUsername,
		s.Config.SMTPPassword,
		s.Config.SMTPHost,
	)

	// Compose email
	to := []string{data.To}
	msg := []byte(fmt.Sprintf(
		"From: %s <%s>\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n\r\n"+
			"%s",
		s.Config.SMTPFromName,
		s.Config.SMTPFrom,
		data.To,
		data.Subject,
		body.String(),
	))

	// Send email
	addr := fmt.Sprintf("%s:%s", s.Config.SMTPHost, s.Config.SMTPPort)
	err = smtp.SendMail(addr, auth, s.Config.SMTPFrom, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendVerificationEmail sends an email verification email
func (s *EmailService) SendVerificationEmail(email string, token string) error {
	verificationURL := fmt.Sprintf("%s/verify-email?token=%s&email=%s",
		s.Config.BaseURL, token, email)

	return s.SendEmail(EmailData{
		To:       email,
		Subject:  "Verify Your Email Address",
		Template: "verification.html",
		Data: map[string]interface{}{
			"VerificationURL": verificationURL,
			"Email":           email,
			"AppName":         "Own PayNet",
		},
	})
}

// SendPasswordResetEmail sends a password reset email
func (s *EmailService) SendPasswordResetEmail(email string, token string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s&email=%s",
		s.Config.BaseURL, token, email)

	return s.SendEmail(EmailData{
		To:       email,
		Subject:  "Reset Your Password",
		Template: "password_reset.html",
		Data: map[string]interface{}{
			"ResetURL": resetURL,
			"Email":    email,
			"AppName":  "Own PayNet",
		},
	})
}

// SendWelcomeEmail sends a welcome email after successful registration
func (s *EmailService) SendWelcomeEmail(email string) error {
	return s.SendEmail(EmailData{
		To:       email,
		Subject:  "Welcome to Own PayNet",
		Template: "welcome.html",
		Data: map[string]interface{}{
			"Email":    email,
			"AppName":  "Own PayNet",
			"LoginURL": fmt.Sprintf("%s/login", s.Config.BaseURL),
		},
	})
}
