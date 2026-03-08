package errors

import (
	"context"
	"net/http"
)

type contextKeyT string

const (
	httpErrorContextKey = contextKeyT("http_error")
)

// HTTPErrorFromContext returns the HTTP error set by a handler via SetHTTPError, or nil.
func HTTPErrorFromContext(ctx context.Context) *httpErrorT {
	if p, ok := ctx.Value(httpErrorContextKey).(**httpErrorT); ok && p != nil {
		return *p
	}
	return nil
}

// SetHTTPError sets the error in the request context if none is set yet. The first error set wins (first failure point in the chain).
func SetHTTPError(ctx context.Context, err *httpErrorT) {
	if err == nil {
		return
	}
	if p, ok := ctx.Value(httpErrorContextKey).(**httpErrorT); ok && p != nil && *p == nil {
		*p = err
	}
}

// responseWriter wraps http.ResponseWriter to track whether the handler already sent headers,
// so the middleware can avoid double-write when sending an error response.
type responseWriter struct {
	http.ResponseWriter
	headerSent bool
}

var _ http.ResponseWriter = (*responseWriter)(nil)

func (w *responseWriter) WriteHeader(code int) {
	w.headerSent = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.headerSent = true
	return w.ResponseWriter.Write(b)
}
