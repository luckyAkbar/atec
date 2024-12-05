package common

import (
	"context"

	sendinblue "github.com/sendinblue/APIv3-go-library/v2/lib"
)

func NewBrevoClient(apiKey string) *sendinblue.APIClient {
	brevoCfg := sendinblue.NewConfiguration()
	brevoCfg.AddDefaultHeader("api-key", apiKey)

	return sendinblue.NewAPIClient(brevoCfg)
}

type Mailer struct {
	senderName  string
	senderEmail string
	brevoClient *sendinblue.APIClient
}

func NewMailer(senderName, senderEmail string, brevoClient *sendinblue.APIClient) *Mailer {
	return &Mailer{
		senderName:  senderName,
		senderEmail: senderEmail,
		brevoClient: brevoClient,
	}
}

type SendEmailInput struct {
	ReceiverName  string
	ReceiverEmail string
	HtmlContent   string
	Subject       string
}

func (m *Mailer) SendEmail(ctx context.Context, input SendEmailInput) (*sendinblue.CreateSmtpEmail, error) {
	body := sendinblue.SendSmtpEmail{
		Sender: &sendinblue.SendSmtpEmailSender{
			Name:  m.senderName,
			Email: m.senderEmail,
		},
		To: []sendinblue.SendSmtpEmailTo{
			{
				Email: input.ReceiverEmail,
				Name:  input.ReceiverName,
			},
		},
		HtmlContent: input.HtmlContent,
		Subject:     input.Subject,
	}

	email, res, err := m.brevoClient.TransactionalEmailsApi.SendTransacEmail(ctx, body)
	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	return &email, nil
}

func (m *Mailer) GetClientName() string {
	return "brevo"
}
