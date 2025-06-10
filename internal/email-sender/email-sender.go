package emailsender

import (
	"fmt"
	"log"
	"sync"

	"github.com/go-mail/mail/v2"
)

type Mailer struct {
	dialer    *mail.Dialer
	sender    string
	queue     chan Message
	waitGroup sync.WaitGroup
}

type Message struct {
	To      string
	Subject string
	Body    string
}

func New(host string, port int, username, password, sender string, workers int) *Mailer {
	dialer := mail.NewDialer(host, port, username, password)

	m := &Mailer{
		dialer: dialer,
		sender: sender,
		queue:  make(chan Message, 100), // Буфер на 100 писем
	}

	// Запускаем N воркеров для обработки очереди
	for i := 0; i < workers; i++ {
		go m.worker()
	}

	return m
}

func (m *Mailer) worker() {
	for msg := range m.queue {
		if err := m.send(msg.To, msg.Subject, msg.Body); err != nil {
			log.Printf("Failed to send email to %s: %v", msg.To, err)
		}
		m.waitGroup.Done()
	}
}

func (m *Mailer) send(to, subject, body string) error {
	msg := mail.NewMessage()
	msg.SetHeader("From", m.sender)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)

	if err := m.dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

// SendAsync добавляет письмо в очередь и возвращает управление сразу
func (m *Mailer) SendAsync(to, subject, body string) {
	m.waitGroup.Add(1)
	m.queue <- Message{To: to, Subject: subject, Body: body}
}

// Close ожидает завершения всех отправок и закрывает канал
func (m *Mailer) Close() {
	m.waitGroup.Wait()
	close(m.queue)
}
