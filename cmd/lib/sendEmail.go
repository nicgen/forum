package lib

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(to string, subject string, body string) error {
	// Email configuration
	from := os.Getenv("GOOGLE_APP_EMAIL")
	password := os.Getenv("GOOGLE_APP_PASSWORD")
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message composition
	message := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", to, subject, body))

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
