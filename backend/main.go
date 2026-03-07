package main

import (
	"log/slog"
	"net/http"
	"os"

	"bitbucket.org/atlassian/unsw-comp3900-app/internal/errors"
	"bitbucket.org/atlassian/unsw-comp3900-app/internal/handlers"
	"bitbucket.org/atlassian/unsw-comp3900-app/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	r := chi.NewRouter()
	r.Use(
		middleware.Logger(log),
		errors.HandleHTTPError,
	)

	r.Group(func(r chi.Router) {
		// common middleware for all JSON Response routes
		r.Use(
			middleware.ContentTypeJSON,
			middleware.RequireAcceptJSON,
		)
		r.Get("/health", handlers.Health)
		r.Get("/error", handlers.Error)
	})

	http.ListenAndServe(":8080", r)
}
