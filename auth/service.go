package auth

import "context"

type Transactor interface {
	Run(ctx context.Context, tx func(Stores) error) error
}

type Service struct{ transactor Transactor }

func NewService(transactor Transactor) *Service {
	return &Service{
		transactor: transactor,
	}
}

type RegisterUser struct {
	Name     string
	Email    string
	Password string
}

func (s *Service) RegisterUser(ctx context.Context, form RegisterUser) (*User, error) {
	return NewUser(form.Name, form.Email, form.Password)
}
