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
		<!DOCTYPE html>
		<html>
		<body style="margin:0; padding:0; background:#f5f5f3; font-family:Arial,sans-serif;">
		<table width="100%%" cellpadding="0" cellspacing="0" style="padding:40px 20px;">
			<tr>
			<td align="center">
				<table width="480" cellpadding="0" cellspacing="0"
					style="background:#ffffff; border-radius:12px; border:1px solid #e8e8e5; overflow:hidden;">
				<tr>
					<td style="padding:40px 40px 32px;">
					<p style="margin:0 0 32px; font-size:12px; font-weight:600; letter-spacing:2px;
								color:#a0a09a; text-transform:uppercase;">Qlass</p>
					<p style="margin:0 0 8px; font-size:22px; font-weight:500; color:#1a1a18; line-height:1.3;">
						Your verification code</p>
					<p style="margin:0 0 32px; font-size:15px; color:#6b6b65; line-height:1.6;">
						Use the code below to continue. It expires in 5 minutes.</p>
					<div style="background:#f5f5f3; border-radius:8px; padding:20px; text-align:center; margin:0 0 32px;">
						<span style="font-size:32px; font-weight:500; letter-spacing:10px; color:#1a1a18;">%s</span>
					</div>
					<p style="margin:0; font-size:13px; color:#a0a09a; line-height:1.6;">
						If you didn't request this, you can safely ignore this email.</p>
					</td>
				</tr>
				<tr>
					<td style="border-top:1px solid #e8e8e5; padding:20px 40px;
							display:flex; justify-content:space-between;">
					<span style="font-size:12px; color:#a0a09a;">© 2025 Qlass</span>
					&nbsp;&nbsp;&nbsp;&nbsp;
					<span style="font-size:12px; color:#a0a09a;">Do not reply to this email</span>
					</td>
				</tr>
				</table>
			</td>
			</tr>
		</table>
		</body>
		</html>`, otp)

	return s.SendEmail([]string{toEmail}, subject, body)
}
