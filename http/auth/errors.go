package auth

import (
	"errors"
	"net/http"

	"github.com/purpose-robot/planet-express/auth"
	"github.com/purpose-robot/planet-express/internal/httputil"
	"github.com/purpose-robot/planet-express/internal/validator"
)

func mapDomainError(err error) error {
	switch {
	case errors.Is(err, auth.ErrRecordNotFound):
		err = httputil.NewPublicError(err, http.StatusNotFound, auth.ErrRecordNotFound.Error())

	case errors.Is(err, auth.ErrDuplicateEmail):
		v := validator.Validator{}
		v.AddFieldError("email", "Email address is already in use")
		err = httputil.NewPublicError(err, http.StatusUnprocessableEntity, auth.ErrDuplicateEmail.Error())
	}

	return err
}
