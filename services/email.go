package services

import (
	"errors"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// client

type MailClient interface {
	Send(message *mail.SGMailV3) (*rest.Response, error)
}

type sendGridMailClient struct {
	sgClient *sendgrid.Client
}

func NewSendGridMailClient() sendGridMailClient {
	return sendGridMailClient{sgClient: sendgrid.NewSendClient(viper.GetString("sendgrid_api_key"))}
}

func (c sendGridMailClient) Send(message *mail.SGMailV3) (*rest.Response, error) {
	return c.sgClient.Send(message)
}

// service

type EmailService struct {
	client    MailClient
	fromEmail *mail.Email
}

func NewEmailService(mailClient MailClient) EmailService {
	return EmailService{
		client:    mailClient,
		fromEmail: mail.NewEmail(viper.GetString("sendgrid_from_name"), viper.GetString("sendgrid_from_email")),
	}
}

func (s EmailService) SendEmail(subject string, content string, to string) error {
	message := mail.NewSingleEmailPlainText(s.fromEmail, subject, mail.NewEmail("", to), content)
	r, err := s.client.Send(message)

	if err != nil {
		return err
	} else if r.StatusCode < 200 || r.StatusCode > 299 {
		zap.S().Error("email failed to send", r.Body)
		return errors.New("failed to send email")
	}

	return nil
}
