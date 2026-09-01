package auth

import (
	"context"
	"errors"
	"fmt"
	"time"
	"uuid"

	"github.com/alexedwards/argon2id"
	"github.com/purpose-robot/planet-express/internal/validator"
)

var (
	ErrRecordNotFound = errors.New("the requested resource could not be found")
	ErrDuplicateEmail = errors.New("user with this email address already exists")
)

type User struct {
	ID        uuid.UUID
	Version   int
	IsActive  bool
	Name      string
	Email     string
	Password  password
	CreatedAt time.Time
}

type password struct {
	Hash      []byte
	Plaintext *string
}

func (p *password) Set(plaintext string) error {
	hash, err := argon2id.CreateHash(plaintext, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	p.Hash = []byte(hash)
	p.Plaintext = &plaintext

	return nil
}

func NewUser(name, email, password string) (*User, error) {
	v := new(validator.Validator)

	v.CheckField(validator.NotBlank(name), "name", "This field cannot be blank")
	v.CheckField(validator.NotBlank(email), "email", "This field cannot be blank")
	v.CheckField(validator.NotBlank(password), "password", "This field cannot be blank")
	v.CheckField(validator.MinRunes(password, 24), "password", "This field must be at least 24 characters")
	v.CheckField(validator.MaxRunes(password, 96), "password", "This field must not be more than 96 characters")

	if !v.Valid() {
		return nil, v
	}

	user := new(User)

	user.ID = uuid.New()
	user.Version = 1
	user.IsActive = false
	user.Name = name
	user.Email = email
	user.CreatedAt = time.Now().UTC()

	if err := user.Password.Set(password); err != nil {
		return nil, fmt.Errorf("create hash from plain-text password for %s: %v", name, err)
	}

	return user, nil
}

type Token struct {
	Hash      []byte
	Plaintext string
	Scope     Scope
	UserID    uuid.UUID
	ExpiresAt time.Time
}

type Scope string

const (
	ScopeActivation     Scope = "activation"
	ScopePasswordReset  Scope = "password_reset"
	ScopeAuthentication Scope = "authentication"
)

type Email struct {
	Name      string
	Plaintext string
	Recipient string
}

type Stores struct {
	Users    UserStore
	Tokens   TokenStore
	Enqueuer EnqueuerStore
}

type UserStore interface {
	Create(ctx context.Context, user *User) error
}

type TokenStore interface {
	Create(ctx context.Context, token *Token) error
}

type EnqueuerStore interface {
	EnqueueActivationEmail(ctx context.Context, email Email) error
	EnqueueResetPasswordEmail(ctx context.Context, email Email) error
}
