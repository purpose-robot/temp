package auth

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/purpose-robot/planet-express/auth"
)

type Transactor struct{ conn *sql.DB }

func NewTransactor(conn *sql.DB) *Transactor {
	return &Transactor{
		conn: conn,
	}
}

func (t *Transactor) Run(ctx context.Context, fn func(tx auth.Stores) error) error {
	tx, err := t.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	stores := auth.Stores{
		Users:  NewUserStore(tx),
		Tokens: NewTokenStore(tx),
	}

	if err := fn(stores); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
