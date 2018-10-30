package syslogalert

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type BufferedMailer struct {
	// chan to receive on
	MessageChan chan Message
	// buffer of Message to store
	Buffer  []Message
	Timeout time.Duration
	Mailer  SendMailer
}

func NewBufferedMailer(timeout time.Duration, mailer SendMailer) *BufferedMailer {
	return &BufferedMailer{
		MessageChan: make(chan Message),
		Buffer:      []Message{},
		Timeout:     10,
		Mailer:      mailer,
	}
}

func (b *BufferedMailer) Start() {
	for {
		// Get first message
		first := <-b.MessageChan
		b.Buffer = append(b.Buffer, first)
		// Push subsequent messages to buffer

	BufferLoop:
		for {
			select {
			case msg := <-b.MessageChan:
				b.Buffer = append(b.Buffer, msg)
			case <-time.After(b.Timeout * time.Second):
				joined := b.JoinMessages()
				header := b.Buffer[0].Header
				b.Mailer.SendMail(header, joined)
				log.Printf("ALERTING: %s %d messages", header, len(b.Buffer))
				break BufferLoop
			}
		}
		// Flush buffer
		b.Buffer = nil
	}
}

func (b *BufferedMailer) JoinMessages() string {
	joinedMessages := make([]string, len(b.Buffer))
	for i, msg := range b.Buffer {
		joinedMessages[i] = fmt.Sprintf("%s\n%s", msg.Header, msg.Body)
	}
	return strings.Join(joinedMessages, "\n")
}

func (b *BufferedMailer) SendMessage(message Message) error {
	log.Printf("QUEUED: %s", message.Header)
	b.MessageChan <- message
	return nil
}
