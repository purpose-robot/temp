package jobs

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/purpose-robot/planet-express/auth"
	"github.com/riverqueue/river"
)

type emailGateway interface {
	Send(ctx context.Context, recipient string, data any, patterns ...string) error
}

type Enqueuer struct {
	tx     *sql.Tx
	client *river.Client[*sql.Tx]
}

func NewEnqueuer(tx *sql.Tx, client *river.Client[*sql.Tx]) *Enqueuer {
	return &Enqueuer{
		tx:     tx,
		client: client,
	}
}

func (e *Enqueuer) EnqueueActivationEmail(ctx context.Context, email auth.Email) error {
	return e.enqueueEmail(ctx, email, "activate-user.tmpl")
}

func (e *Enqueuer) EnqueueResetPasswordEmail(ctx context.Context, email auth.Email) error {
	return e.enqueueEmail(ctx, email, "reset-password.tmpl")
}

func (e *Enqueuer) enqueueEmail(ctx context.Context, email auth.Email, template string) error {
	args := EmailArgs{
		Template:  template,
		Name:      email.Name,
		Plaintext: email.Plaintext,
		Recipient: email.Recipient,
	}

	_, err := e.client.InsertTx(ctx, e.tx, args, nil)
	if err != nil {
		return fmt.Errorf("insert job for %s email: %w", template, err)
	}

	return nil
}
