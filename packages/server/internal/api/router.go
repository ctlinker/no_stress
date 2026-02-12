package api

import (
	"server/internal/api/middlewares"
	"server/internal/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(db *db.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", Health)

	r.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.With(middleware.AllowContentType("application/json")).Post("/register", CreateUserHandler(db))
			r.With(middleware.AllowContentType("application/json")).Post("/connect", UserConnectionHandler(db))
			r.With(middlewares.WithAccessToken).Get("/check", Check)
		})
	})
	return r
}
