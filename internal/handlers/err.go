package handlers

import (
	"net/http"

	httpErrors "bitbucket.org/atlassian/unsw-comp3900-app/internal/errors"
)

func Error(w http.ResponseWriter, r *http.Request) {
	httpErrors.SetHTTPError(r.Context(), httpErrors.MyCustomError())
	// middleware will send the error response if no response was written yet
}
