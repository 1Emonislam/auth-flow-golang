package services

import (
	"own-paynet/config"
	"own-paynet/utils/email"
)

type EmailService struct {
	emailService *email.EmailService
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{
		emailService: email.NewEmailService(cfg),
	}
}

func (s *EmailService) SendOTPEmail(to, otp string) error {
	return s.emailService.SendEmail(email.EmailData{
		To:       to,
		Subject:  "Your 2FA Verification Code",
		Template: "two_factor.html",
		Data: map[string]interface{}{
			"OTP":     otp,
			"AppName": "Manty Pay",
		},
	})
}
