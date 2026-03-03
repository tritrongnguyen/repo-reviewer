package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/tritrongnguyen/repo-reviewer.git/internal/handler"
)

func (s *Server) RegisterRoutes() http.Handler {

	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", handler.Health)

	// r.Route("/api/v1", func(r chi.Router) {
	// 	r.Post("/webhook/github", handler.GithubWebhook)
	// })

	return r
}
