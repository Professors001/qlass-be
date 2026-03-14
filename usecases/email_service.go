package usecases

type EmailService interface {
	SendOTP(toEmail string, otp string) error
	SendEmail(to []string, subject string, body string) error
}
