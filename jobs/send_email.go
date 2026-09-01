package jobs

import (
	"context"
	"fmt"

	"github.com/riverqueue/river"
)

type EmailArgs struct {
	Template  string
	Name      string
	Plaintext string
	Recipient string
}

func (EmailArgs) Kind() string {
	return "emails"
}

func (EmailArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		MaxAttempts: 5,
	}
}

type EmailWorker struct {
	email emailGateway
	river.WorkerDefaults[EmailArgs]
}

func NewEmailWorker(email emailGateway) *EmailWorker {
	return &EmailWorker{
		email: email,
	}
}

func (w *EmailWorker) Work(ctx context.Context, job *river.Job[EmailArgs]) error {
	data := map[string]any{
		"name":      job.Args.Name,
		"plaintext": job.Args.Plaintext,
	}

	err := w.email.Send(ctx, job.Args.Recipient, data, job.Args.Template)
	if err != nil {
		return fmt.Errorf("send %s email to %s: %w", job.Args.Template, job.Args.Recipient, err)
	}

	return nil
}
