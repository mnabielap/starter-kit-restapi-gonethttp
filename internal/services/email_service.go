package services

import (
	"fmt"
	"net/smtp"

	"starter-kit-restapi-gonethttp/config"
)

type EmailService interface {
	SendEmail(to, subject, body string) error
	SendResetPasswordEmail(to, token string) error
	SendVerificationEmail(to, token string) error
}

type emailService struct {
	cfg *config.Config
}

func NewEmailService(cfg *config.Config) EmailService {
	return &emailService{cfg: cfg}
}

func (s *emailService) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.cfg.SMTP.Username, s.cfg.SMTP.Password, s.cfg.SMTP.Host)
	
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/plain; charset=\"utf-8\"\r\n"+
		"\r\n"+
		"%s\r\n", to, subject, body))

	addr := fmt.Sprintf("%s:%d", s.cfg.SMTP.Host, s.cfg.SMTP.Port)
	
	if s.cfg.Env == "test" {
		return nil // Do not send email in test mode
	}

	return smtp.SendMail(addr, auth, s.cfg.SMTP.From, []string{to}, msg)
}

func (s *emailService) SendResetPasswordEmail(to, token string) error {
	subject := "Reset Password"
	// Replace with your frontend URL
	resetURL := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", token)
	text := fmt.Sprintf("Dear user,\n\nTo reset your password, click on this link: %s\n\nIf you did not request any password resets, then ignore this email.", resetURL)
	return s.SendEmail(to, subject, text)
}

func (s *emailService) SendVerificationEmail(to, token string) error {
	subject := "Email Verification"
	// Replace with your frontend URL
	verifyURL := fmt.Sprintf("http://localhost:3000/verify-email?token=%s", token)
	text := fmt.Sprintf("Dear user,\n\nTo verify your email, click on this link: %s\n\nIf you did not create an account, then ignore this email.", verifyURL)
	return s.SendEmail(to, subject, text)
}