package mailer

import (
	"bytes"
	"context"
	"fmt"
	htmlTemplate "html/template"
	textTemplate "text/template"
	"time"

	"github.com/purpose-robot/planet-express/assets"
	"github.com/wneessen/go-mail"
)

type EmailGateway struct {
	from   string
	client *mail.Client
}

func NewEmailGateway(host string, port int, from, username, password string) (*EmailGateway, error) {
	client, err := mail.NewClient(
		host,
		mail.WithPort(port),
		mail.WithTimeout(10*time.Second),
		mail.WithUsername(username),
		mail.WithPassword(password),
		mail.WithSMTPAuth(mail.SMTPAuthLogin),
	)
	if err != nil {
		return nil, err
	}

	return &EmailGateway{from: from, client: client}, nil
}

func (e *EmailGateway) Send(ctx context.Context, recipient string, data any, patterns ...string) error {
	for i := range patterns {
		patterns[i] = "emails/" + patterns[i]
	}

	msg := mail.NewMsg()

	if err := msg.From(e.from); err != nil {
		return fmt.Errorf("set sender in email body: %v", err)
	}

	if err := msg.To(recipient); err != nil {
		return fmt.Errorf("set recipient in email body: %v", err)
	}

	t, err := textTemplate.New("").ParseFS(assets.Emails, patterns...)
	if err != nil {
		return fmt.Errorf("parse text templates %v: %v", patterns, err)
	}

	subject := new(bytes.Buffer)
	if err := t.ExecuteTemplate(subject, "subject", data); err != nil {
		return fmt.Errorf("execute text template for subject field: %v", err)
	}

	msg.Subject(subject.String())

	content := new(bytes.Buffer)
	if err := t.ExecuteTemplate(content, "content", data); err != nil {
		return fmt.Errorf("execute text template for content field: %v", err)
	}

	msg.SetBodyString(mail.TypeTextPlain, content.String())

	if t.Lookup("htmlContent") != nil {
		t, err := htmlTemplate.New("").ParseFS(assets.Emails, patterns...)
		if err != nil {
			return fmt.Errorf("parse html template %v: %v", patterns, err)
		}

		htmlContent := new(bytes.Buffer)
		if err := t.ExecuteTemplate(htmlContent, "htmlContent", data); err != nil {
			return fmt.Errorf("execute html template for content field: %v", err)
		}

		msg.AddAlternativeString(mail.TypeTextHTML, htmlContent.String())
	}

	err = e.client.DialAndSendWithContext(ctx, msg)
	if err != nil {
		return fmt.Errorf("send email to recipient %s: %v", recipient, err)
	}

	return nil
}
