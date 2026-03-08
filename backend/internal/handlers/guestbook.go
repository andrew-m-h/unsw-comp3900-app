package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	httpErrors "bitbucket.org/atlassian/unsw-comp3900-app/internal/errors"
	"bitbucket.org/atlassian/unsw-comp3900-app/internal/guestbook"
)

type createGuestbookRequest struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

// Guestbook returns handlers for the /api/guestbook routes (create and list).
func Guestbook(client *guestbook.Client) (create, list http.HandlerFunc) {
	create = func(w http.ResponseWriter, r *http.Request) {
		var req createGuestbookRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpErrors.SetHTTPError(r.Context(), httpErrors.HTTPError(http.StatusBadRequest, err))
			return
		}
		if req.Name == "" || req.Message == "" {
			httpErrors.SetHTTPError(r.Context(), httpErrors.HTTPError(http.StatusBadRequest, errors.New("name and message are required")))
			return
		}
		entry, err := client.CreateEntry(r.Context(), req.Name, req.Message)
		if err != nil {
			httpErrors.SetHTTPError(r.Context(), httpErrors.HTTPError(http.StatusInternalServerError, err))
			return
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(entry)
	}

	list = func(w http.ResponseWriter, r *http.Request) {
		entries, err := client.ListEntries(r.Context())
		if err != nil {
			httpErrors.SetHTTPError(r.Context(), httpErrors.HTTPError(http.StatusInternalServerError, err))
			return
		}
		if entries == nil {
			entries = []guestbook.Entry{}
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(entries)
	}

	return create, list
}
