package auth

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/purpose-robot/planet-express/auth"
	"github.com/purpose-robot/planet-express/internal/httputil"
)

type Handler struct {
	service  *auth.Service
	renderer *httputil.Renderer
	sessions *scs.SessionManager
}

func NewHandler(service *auth.Service, renderer *httputil.Renderer, sessions *scs.SessionManager) *Handler {
	return &Handler{
		service:  service,
		renderer: renderer,
		sessions: sessions,
	}
}

type signupForm struct {
	Name     string `form:"Name"`
	Email    string `form:"Email"`
	Password string `form:"Password"`
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var form signupForm

	data := httputil.NewTemplateData()
	data.Form = form

	if err := h.renderer.Render(w, http.StatusOK, data, "pages/register-user.tmpl"); err != nil {
		h.renderer.InternalServerError(w, r, err)
	}
}

func (h *Handler) SignupPost(w http.ResponseWriter, r *http.Request) {
	var form signupForm

	if err := httputil.DecodePostForm(w, r, &form); err != nil {
		h.renderer.BadRequest(w, r, err)
		return
	}

	user, err := h.service.RegisterUser(r.Context(), auth.RegisterUser{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	})
	if err != nil {
		mapDomainError(err)
		return
	}

	err = h.sessions.RenewToken(r.Context())
	if err != nil {
		h.renderer.InternalServerError(w, r, err)
		return
	}

	h.sessions.Put(r.Context(), "authenticatedUser", user)

	data := httputil.NewTemplateData()
	data.Flash = "Your account was successfully registered. Please check your inbox"

	if err := h.renderer.Render(w, http.StatusAccepted, data, "page/activate-user.tmpl"); err != nil {
		h.renderer.InternalServerError(w, r, err)
	}
}
