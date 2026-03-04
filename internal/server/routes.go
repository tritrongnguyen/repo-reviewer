package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/tritrongnguyen/repo-reviewer.git/internal/handler"
	"github.com/tritrongnguyen/repo-reviewer.git/internal/repository"
	"github.com/tritrongnguyen/repo-reviewer.git/internal/service"
)

func (s *Server) RegisterRoutes() http.Handler {

	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	// Cors setting
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Inject dependencies
	db := s.db.GetDB()

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	authService := service.NewAuthService(userRepo, sessionRepo)

	authHandler := handler.NewAuthHandler(authService)

	// Routes
	r.Get("/health", handler.Health)

	r.Route("/api/v1", func(r chi.Router) {

		r.Route("/auth", func(r chi.Router) {
			r.Post("/signup", authHandler.SignUp)
			r.Post("/login", authHandler.Login)
			r.Post("/reset-password", authHandler.ResetPassword)
		})

		r.Post("/webhook/github", handler.GithubWebhook)
	})

	return r
}
