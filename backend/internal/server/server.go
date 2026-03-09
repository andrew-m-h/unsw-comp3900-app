package server

import (
	"log/slog"
	"net/http"

	"bitbucket.org/atlassian/unsw-comp3900-app/internal/errors"
	"bitbucket.org/atlassian/unsw-comp3900-app/internal/guestbook"
	"bitbucket.org/atlassian/unsw-comp3900-app/internal/handlers"
	"bitbucket.org/atlassian/unsw-comp3900-app/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// NewHandler builds the same HTTP handler used in production (Chi router with health, API, guestbook).
// Callers can use it with http.ListenAndServe or httptest.NewServer for integration tests.
func NewHandler(gbClient *guestbook.Client, log *slog.Logger, version string) http.Handler {
	guestbookCreate, guestbookList := handlers.Guestbook(gbClient)

	r := chi.NewRouter()
	r.Use(
		middleware.Logger(log),
		errors.HandleHTTPError,
	)

	r.Group(func(r chi.Router) {
		r.Use(
			middleware.ResponseContentTypeJSON,
			middleware.RequireAcceptJSON,
		)
		r.Get("/health", handlers.Health(version))
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(
			middleware.ResponseContentTypeJSON,
			middleware.RequireAcceptJSON,
			middleware.RequireContentTypeJSONForBody,
		)

		r.Get("/health", handlers.Health(version))
		r.Route("/guestbook", func(r chi.Router) {
			r.Get("/", guestbookList)
			r.Post("/", guestbookCreate)
			r.Delete("/{id}", handlers.DeleteGuestbookEntry(gbClient))
			r.Get("/{id}", handlers.GuestbookGet(gbClient))
		})
	})

	return r
}
