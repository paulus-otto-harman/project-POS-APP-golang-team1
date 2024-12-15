package service

import (
	"bytes"
	"context"
	"github.com/mailersend/mailersend-go"
	"go.uber.org/zap"
	"html/template"
	"project/config"
	"time"
)

type EmailService interface {
	Send(to, subject, template string, data interface{}) (string, error)
}

type emailService struct {
	log    *zap.Logger
	Mailer *mailersend.Mailersend
	sender mailersend.From
}

func NewEmailService(config config.EmailConfig, log *zap.Logger) EmailService {
	ms := mailersend.NewMailersend(config.ApiKey)
	sender := mailersend.From{
		Name:  config.FromName,
		Email: config.FromEmail,
	}
	return &emailService{Mailer: ms, log: log, sender: sender}
}

func (s *emailService) Send(to, subject, htmlTemplate string, data interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	message := s.Mailer.Email.NewMessage()

	recipient := mailersend.Recipient{
		Email: to,
	}
	message.SetFrom(s.sender)
	message.SetRecipients([]mailersend.Recipient{recipient})
	message.SetSubject(subject)

	tmpl, err := template.ParseFiles("../email/" + htmlTemplate + ".html")
	if err != nil {
		return "", err
	}

	// Apply template dengan data
	var body bytes.Buffer
	if err = tmpl.Execute(&body, data); err != nil {
		return "", err
	}
	s.log.Info("Sending email", zap.String("to", to), zap.String("subject", subject))
	s.log.Info("Sending email", zap.String("body", body.String()))
	message.SetHTML(body.String())

	var response *mailersend.Response
	response, err = s.Mailer.Email.Send(ctx, message)
	if err != nil {
		return "", err
	}

	return response.Header.Get("X-Message-Id"), nil
}
