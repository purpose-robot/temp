package http

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/purpose-robot/planet-express/assets"
	"github.com/purpose-robot/planet-express/auth"
	"github.com/purpose-robot/planet-express/internal/httputil"

	httpauth "github.com/purpose-robot/planet-express/http/auth"
)

type Server struct {
	authSvc  *auth.Service
	renderer *httputil.Renderer
	sessions *scs.SessionManager
}

func NewServer(authSvc *auth.Service, renderer *httputil.Renderer, sessions *scs.SessionManager) *Server {
	return &Server{
		authSvc:  authSvc,
		renderer: renderer,
		sessions: sessions,
	}
}

func (s *Server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (s *Server) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (s *Server) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (s *Server) Handle() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(assets.Static))

	// mux.Handle("GET /{$}", nil)

	authHandler := httpauth.NewHandler(s.authSvc, s.renderer, s.sessions)

	mux.HandleFunc("GET /sign-up", authHandler.Signup)
	mux.HandleFunc("POST /sign-up", authHandler.SignupPost)

	mux.HandleFunc("GET /", s.renderer.NotFound)

	return s.logRequest(s.recoverPanic(s.securityHeaders(s.sessions.LoadAndSave(mux))))
}
