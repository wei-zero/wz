package sgmailer

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/wei-zero/wz/mailer"
)

type SendgridMailer struct {
	apiKey string
}

func NewSendgridMailer(apiKey string) *SendgridMailer {
	return &SendgridMailer{apiKey: apiKey}
}

// Ref: https://github.com/sendgrid/sendgrid-go/blob/main/use-cases/transactional-templates-with-mailer-helper.md
func (s *SendgridMailer) SendByTemplate(from *mailer.Email, to *mailer.Email, subject string, templateID string, data map[string]string) error {
	personalization := mail.NewPersonalization()
	personalization.AddTos(mail.NewEmail(to.Name, to.Address))
	for k, v := range data {
		personalization.SetDynamicTemplateData(k, v)
	}

	m := mail.NewV3Mail()
	m.Subject = subject
	m.SetFrom(mail.NewEmail(from.Name, from.Address))
	m.SetTemplateID(templateID)
	m.AddPersonalizations(personalization)

	request := sendgrid.GetRequest(s.apiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	var Body = mail.GetRequestBody(m)
	request.Body = Body
	res, err := sendgrid.API(request)
	if err != nil {
		return err
	}
	if res.StatusCode != 202 {
		return fmt.Errorf("send mail failed, code=%d, body=%s", res.StatusCode, res.Body)
	}

	return nil
}
