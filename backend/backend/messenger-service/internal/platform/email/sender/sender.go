package sender

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/platform/email/config"
	"gopkg.in/gomail.v2"
)

type Sender struct {
	dialer *gomail.Dialer
	from string
}

func NewSender(cfg *config.EmailConfig) *Sender {
	d := gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass)
	return &Sender{dialer: d, from: cfg.SMTPFrom}
}

func (s *Sender) SendVerificationCode(toEmail, code string, purpose string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Verification Code - Messenger")

	var title, description string
	switch purpose {
	case "VERIFY_EMAIL":
		title = "Verify your email"
		description = "Please use the following code to verify your email"
	case "CHANGE_PASSWORD":
		title = "Change your password"
		description = "Please use the following code to change your password"
	default:
		return errors.New("Unknown email purpose")
	}

	htmlBody := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; padding: 20px; background-color: #f4f4f4;">
			<div style="background-color: #ffffff; padding: 20px; border-radius: 8px; max-width: 500px; margin: auto;">
				<h2 style="color: #333;">Welcome! %s</h2>
				<p>%s</p>
				<h1 style="color: #007bff; letter-spacing: 5px; text-align: center;">%s</h1>
				<p style="color: #777; font-size: 12px;">This code expires in 15 minutes.</p>
			</div>
		</div>
	`, title, description, code)

	m.SetBody("text/html", htmlBody)

	if err := s.dialer.DialAndSend(m); err != nil {
		log.Error().Err(err).Str("to", toEmail).Str("purpose", purpose).Msg("Failed to send verification email")
		return err
	}

	log.Debug().Str("to", toEmail).Str("purpose", purpose).Msg("Verification email sent")
	return nil
}