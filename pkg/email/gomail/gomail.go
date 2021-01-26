package gomail

import (
	"github.com/go-gomail/gomail"
)

func SendEmail(host string, port int, userName, password, from, fromName string, to []string, subject,
	content string, asHTML bool) error {
	emailMsg := gomail.NewMessage()
	emailMsg.SetAddressHeader("From", from, fromName)

	emailMsg.SetHeader("To", to...)

	emailMsg.SetHeader("Subject", subject)

	if asHTML {
		emailMsg.SetBody("text/html", content)
	} else {
		emailMsg.SetBody("text", content)
	}

	return gomail.NewDialer(host, port, userName, password).DialAndSend(emailMsg)
}
