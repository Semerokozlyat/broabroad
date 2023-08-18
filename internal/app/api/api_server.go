package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"broabroad/internal/app/config"
)

type Server struct {
	*http.Server
}

func NewServer(cfg *config.Config, seekReqCreator SeekRequestsCreator) *Server {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.BasicAuth("bro-abroad", map[string]string{
		cfg.Access.User: cfg.Access.Password,
	}))
	r.Use(middleware.Recoverer)

	r.Method(http.MethodGet, "/", newRootHandler())
	r.Method(http.MethodGet, "/members", newGetMembersHandler())
	r.Method(http.MethodPost, "/members", newAddMemberHandler())
	r.Method(http.MethodPost, "/seek_requests", newCreateSeekRequestHandler(seekReqCreator))

	server := &Server{
		Server: &http.Server{
			Addr:        cfg.HTTPServer.Address,
			Handler:     r,
			ReadTimeout: cfg.HTTPServer.Timeout,
			IdleTimeout: cfg.HTTPServer.IdleTimeout,
		},
	}
	return server
}

func (s *Server) Run() error {
	return s.Server.ListenAndServe()
}
