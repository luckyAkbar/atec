package common

import (
	"context"

	sendinblue "github.com/sendinblue/APIv3-go-library/v2/lib"
	"github.com/sweet-go/stdlib/helper"
)

// NewBrevoClient create a new brevo client instance
func NewBrevoClient(apiKey string) *sendinblue.APIClient {
	brevoCfg := sendinblue.NewConfiguration()
	brevoCfg.AddDefaultHeader("api-key", apiKey)

	return sendinblue.NewAPIClient(brevoCfg)
}

// Mailer contains utility relates to mailing functionalities
type Mailer struct {
	senderName  string
	senderEmail string
	brevoClient *sendinblue.APIClient
}

// NewMailer creates a new mailer instance
func NewMailer(senderName, senderEmail string, brevoClient *sendinblue.APIClient) *Mailer {
	return &Mailer{
		senderName:  senderName,
		senderEmail: senderEmail,
		brevoClient: brevoClient,
	}
}

// SendEmailInput input to send email
type SendEmailInput struct {
	ReceiverName  string
	ReceiverEmail string
	HTMLContent   string
	Subject       string
}

// SendEmail simple email sending functionality utilizing brevo
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
		HtmlContent: input.HTMLContent,
		Subject:     input.Subject,
	}

	email, res, err := m.brevoClient.TransactionalEmailsApi.SendTransacEmail(ctx, body)
	defer helper.WrapCloser(res.Body.Close)

	if err != nil {
		return nil, err
	}

	return &email, nil
}
