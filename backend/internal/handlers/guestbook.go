package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	httpErrors "bitbucket.org/atlassian/unsw-comp3900-app/internal/errors"
	"bitbucket.org/atlassian/unsw-comp3900-app/internal/guestbook"
	"github.com/go-chi/chi/v5"
)

type CreateGuestbookRequest struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

// Guestbook returns handlers for the /api/guestbook routes (create and list).
func Guestbook(client *guestbook.Client) (create, list http.HandlerFunc) {
	create = func(w http.ResponseWriter, r *http.Request) {
		var req CreateGuestbookRequest
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

// GuestbookGet returns a single guestbook entry by id. GET /api/guestbook/{id} → 200 with entry or 404.
func GuestbookGet(client *guestbook.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			httpErrors.SetHTTPError(r.Context(), httpErrors.HTTPError(http.StatusBadRequest, errors.New("id is required")))
			return
		}
		entry, err := client.GetEntry(r.Context(), id)
		if err != nil {
			if errors.Is(err, guestbook.ErrNotFound) {
				httpErrors.SetHTTPError(r.Context(), httpErrors.HTTPError(http.StatusNotFound, err))
				return
			}
			httpErrors.SetHTTPError(r.Context(), httpErrors.HTTPError(http.StatusInternalServerError, err))
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(entry)
	}
}

// DeleteGuestbookEntry deletes a guestbook entry by id. DELETE /api/guestbook/{id} → 204 or 404 if not found.
func DeleteGuestbookEntry(client *guestbook.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			httpErrors.SetHTTPError(r.Context(), httpErrors.HTTPError(http.StatusBadRequest, errors.New("id is required")))
			return
		}
		err := client.DeleteEntry(r.Context(), id)
		if err != nil {
			if errors.Is(err, guestbook.ErrNotFound) {
				httpErrors.SetHTTPError(r.Context(), httpErrors.HTTPError(http.StatusNotFound, err))
				return
			}
			httpErrors.SetHTTPError(r.Context(), httpErrors.HTTPError(http.StatusInternalServerError, err))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
