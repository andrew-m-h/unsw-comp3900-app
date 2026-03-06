package errors

import (
	"context"
	"encoding/json"
	"net/http"
)

func HandleHTTPError(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err *httpErrorT
		r = r.WithContext(context.WithValue(r.Context(), httpErrorContextKey, &err))

		ww := &responseWriter{ResponseWriter: w}
		next.ServeHTTP(ww, r)

		httpError := HTTPErrorFromContext(r.Context())
		if httpError != nil && !ww.headerSent {
			writeJSONErrorResponse(w, httpError.Code, httpError.Err.Error())
		}
	})
}

// errorResponse is the JSON shape for error responses.
type errorResponse struct {
	Error string `json:"error"`
}

// WriteErrorResponse sets Content-Type: application/json, writes the status code, and sends a JSON body {"error": "message"}.
func writeJSONErrorResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: message})
}
