package syslogalert

import (
	"encoding/json"
	"os"

	"github.com/go-gomail/gomail"
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
	gomailer.SetHeader("Subject", header)
	gomailer.SetBody("text/plain", body)

	d := gomail.NewDialer(m.SMTPServer, m.SMTPPort, m.Username, m.Password)
	return d.DialAndSend(gomailer)
}
