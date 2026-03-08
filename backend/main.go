package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"bitbucket.org/atlassian/unsw-comp3900-app/internal/errors"
	"bitbucket.org/atlassian/unsw-comp3900-app/internal/guestbook"
	"bitbucket.org/atlassian/unsw-comp3900-app/internal/handlers"
	"bitbucket.org/atlassian/unsw-comp3900-app/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	gbClient, err := guestbook.NewClient(context.Background())
	if err != nil {
		log.Error("guestbook client", "err", err)
		os.Exit(1)
	}
	guestbookCreate, guestbookList := handlers.Guestbook(gbClient)

	r := chi.NewRouter()
	r.Use(
		// request logging
		// add logger to context
		middleware.Logger(log),
		// handle HTTP errors raised during request processing
		errors.HandleHTTPError,
	)

	// expose a /health check endpoint
	r.Group(func(r chi.Router) {
		r.Use(
			// respond with Content-Type: application/json
			middleware.ResponseContentTypeJSON,
			// require Accept: application/json
			middleware.RequireAcceptJSON,
		)
		r.Get("/health", handlers.Health(Version))
	})

	// API routes under /api
	r.Route("/api", func(r chi.Router) {
		r.Use(
			// respond with Content-Type: application/json
			middleware.ResponseContentTypeJSON,
			// require Accept: application/json
			middleware.RequireAcceptJSON,
			// require Content-Type: application/json for POST, PUT, PATCH requests
			middleware.RequireContentTypeJSONForBody,
		)

		r.Get("/error", handlers.Error)
		r.Get("/health", handlers.Health(Version))
		r.Get("/guestbook/{id}", handlers.GuestbookGet(gbClient)) // before /guestbook so GET /api/guestbook/{id} matches
		r.Route("/guestbook", func(r chi.Router) {
			r.Get("/", guestbookList)
			r.Post("/", guestbookCreate)
		})
	})

	http.ListenAndServe(":8080", r)
}
