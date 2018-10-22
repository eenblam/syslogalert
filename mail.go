package syslogalert

import (
	"encoding/json"
	"os"

	"gopkg.in/gomail.v2"
)

type SendMailer interface {
	SendMail(header, body string) error
}

type Mailer struct {
	From       string
	To         []string
	SMTPServer string
	SMTPPort   int
	Username   string
	Password   string
}

func NewMailer(filename string) (*Mailer, error) {
	mailer := &Mailer{}
	ruleFile, readErr := os.Open(filename)
	if readErr != nil {
		return nil, readErr
	}
	jsonParser := json.NewDecoder(ruleFile)
	parseErr := jsonParser.Decode(mailer)
	if parseErr != nil {
		return nil, parseErr
	}
	return mailer, nil
}

func (m *Mailer) SendMail(header, body string) error {
	gomailer := gomail.NewMessage()
	gomailer.SetHeader("From", m.From)
	gomailer.SetHeader("To", m.To...)
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	gomailer.SetHeader("Subject", header)
	gomailer.SetBody("text/plain", body)
	//m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer(m.SMTPServer, m.SMTPPort, m.Username, m.Password)
	return d.DialAndSend(gomailer)
}
