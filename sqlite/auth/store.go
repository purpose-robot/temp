package auth

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/purpose-robot/planet-express/auth"
	"github.com/purpose-robot/planet-express/sqlite"
)

type dbtx interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type UserStore struct{ db dbtx }

// Compile-time interface guard.
var _ auth.UserStore = (*UserStore)(nil)

func NewUserStore(db dbtx) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (s *UserStore) Create(ctx context.Context, user *auth.User) error {
	query := `
		INSERT INTO users (id, version, is_active, name, email, password_hash, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	args := []any{
		user.ID,
		user.Version,
		user.IsActive,
		user.Name,
		user.Email,
		user.Password.Hash,
		user.CreatedAt,
	}

	_, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		if !sqlite.CheckConstraint(err, "users_email") {
			return fmt.Errorf("insert user %s: %v", user.Email, err)
		}

		return fmt.Errorf("user %s: %w", user.Email, auth.ErrDuplicateEmail)
	}

	return nil
}

type TokenStore struct{ db dbtx }

// Compile-time interface guard.
var _ auth.TokenStore = (*TokenStore)(nil)

func NewTokenStore(db dbtx) *TokenStore {
	return &TokenStore{
		db: db,
	}
}

func (s *TokenStore) Create(ctx context.Context, token *auth.Token) error {
	query := `
		INSERT INTO tokens (hash, scope, user_id, expires_at)
		VALUES (?, ?, ?, ?)`

	args := []any{
		token.Hash,
		token.Scope,
		token.UserID,
		token.ExpiresAt,
	}

	_, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("insert %s token for user %s: %w", token.Scope, token.UserID, err)
	}

	return nil
}
