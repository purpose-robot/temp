package httputil

import "net/http"

type PublicError struct {
	err     error
	Status  int
	Message string
}

func (e *PublicError) Error() string {
	return e.Message
}

func NewPublicError(err error, status int, message string) *PublicError {
	return &PublicError{
		err:     err,
		Status:  status,
		Message: message,
	}
}

func (h *Renderer) Error(w http.ResponseWriter, r *http.Request, err error) {}

func (h *Renderer) NotFound(w http.ResponseWriter, r *http.Request) {
	h.Error(w, r, NewPublicError(nil, http.StatusNotFound, "the requested resource could not be found"))
}

func (h *Renderer) BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	h.Error(w, r, NewPublicError(err, http.StatusBadRequest, "the server could not understand your request"))
}

func (h *Renderer) InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	h.Error(w, r, NewPublicError(err, http.StatusInternalServerError, "the server encountered a problem and could not process your request"))
}
