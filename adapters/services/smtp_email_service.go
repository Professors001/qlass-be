package services

import (
	"fmt"
	"net/smtp"
	"qlass-be/usecases"
)

type smtpEmailService struct {
	host     string
	port     string
	user     string
	password string
}

func NewSMTPEmailService(host, port, user, password string) usecases.EmailService {
	return &smtpEmailService{
		host:     host,
		port:     port,
		user:     user,
		password: password,
	}
}

func (s *smtpEmailService) SendEmail(to []string, subject string, body string) error {
	auth := smtp.PlainAuth("", s.user, s.password, s.host)

	// Build the email headers and body
	// Using MIME format allows us to send HTML emails easily
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"+
		"%s\r\n", to[0], subject, body))

	address := fmt.Sprintf("%s:%s", s.host, s.port)

	return smtp.SendMail(address, auth, s.user, to, msg)
}

func (s *smtpEmailService) SendOTP(toEmail string, otp string) error {
	subject := "Your Qlass Verification Code"

	// Simple HTML template for the OTP
	body := fmt.Sprintf(`
        <div style="font-family: Arial, sans-serif; padding: 20px; text-align: center;">
            <h2>Welcome to Qlass!</h2>
            <p>Your verification code is:</p>
            <h1 style="color: #4A90E2; letter-spacing: 5px;">%s</h1>
            <p>This code will expire in 5 minutes.</p>
        </div>
    `, otp)

	return s.SendEmail([]string{toEmail}, subject, body)
}
