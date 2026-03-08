package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"bitbucket.org/atlassian/unsw-comp3900-app/internal/guestbook"
	"bitbucket.org/atlassian/unsw-comp3900-app/internal/server"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	gbClient, err := guestbook.NewClient(context.Background())
	if err != nil {
		log.Error("guestbook client", "err", err)
		os.Exit(1)
	}

	http.ListenAndServe(":8080", server.NewHandler(gbClient, log, Version))
}
