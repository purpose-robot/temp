package httputil

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/go-playground/form/v4"
)

var decoder = form.NewDecoder()

func DecodePostForm(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)

	err := r.ParseForm()
	if err != nil {
		return err
	}

	return decodeURLValues(r.PostForm, dst)
}

func decodeURLValues(v url.Values, dst any) error {
	err := decoder.Decode(dst, v)
	if err != nil {
		_, ok := errors.AsType[*form.InvalidDecoderError](err)
		if ok {
			panic(err)
		}
	}

	return err
}
