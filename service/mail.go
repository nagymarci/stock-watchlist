package service

import (
	"fmt"
	"net/smtp"
	"os"
)

// smtpServer data to smtp server
type smtpServer struct {
	host string
	port string
}

// Address URI to smtp server
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

type mail struct{}

func NewMail() *mail {
	return &mail{}
}

func (m *mail) SendNotification(profileName string, removed, added, currentStocks []string, email string) error {
	// Sender data.
	from := os.Getenv("SMPT_SENDER_USERNAME")
	password := os.Getenv("SMPT_SENDER_PASSWORD")
	// Receiver email address.
	to := []string{
		email,
	}
	// smtp server configuration.
	smtpServer := smtpServer{host: os.Getenv("SMPT_SERVER_HOST"), port: os.Getenv("SMPT_SERVER_PORT")}
	// Message.
	message := []byte(fmt.Sprintf("To: %v\n"+
		"Subject: %v changed!\n"+
		"\n"+
		"Recommendations in your %v profile has changed.\n\n"+
		"Removed stocks: %+v\n"+
		"Added stocks: %+v\n\n"+
		"Currently recommended stocks: %+v", to[0], profileName, profileName, removed, added, currentStocks))
	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpServer.host)
	// Sending email.
	err := smtp.SendMail(smtpServer.Address(), auth, from, to, message)

	return err
}
