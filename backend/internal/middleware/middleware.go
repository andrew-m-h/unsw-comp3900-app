package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	httpErrors "bitbucket.org/atlassian/unsw-comp3900-app/internal/errors"
)

type contextKeyT string

const (
	loggerContextKey               = contextKeyT("logger")
	responseStatusCodeContextKey   = contextKeyT("response_status_code")
	responseBytesWrittenContextKey = contextKeyT("response_bytes_written")
)

// responseWriter intercepts WriteHeader and Write so status and bytes written are updated in context.
type responseWriter struct {
	http.ResponseWriter
	status  *int
	written *int64
}

// Ensure responseWriter implements http.ResponseWriter
var _ http.ResponseWriter = &responseWriter{}

func (w *responseWriter) WriteHeader(code int) {
	*w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// ResponseStatusFromContext returns the response status code from context (set by Logger middleware after the handler runs). Returns 0 if not set.
func ResponseStatusFromContext(ctx context.Context) int {
	if p, ok := ctx.Value(responseStatusCodeContextKey).(*int); ok && p != nil {
		return *p
	}
	return 0
}

// ResponseBytesWrittenFromContext returns the response body bytes written from context. Returns 0 if not set.
func ResponseBytesWrittenFromContext(ctx context.Context) int64 {
	if p, ok := ctx.Value(responseBytesWrittenContextKey).(*int64); ok && p != nil {
		return *p
	}
	return 0
}

// LoggerFromContext returns the logger injected by Logger middleware, or the default logger if not set.
func LoggerFromContext(ctx context.Context) *slog.Logger {
	if log, ok := ctx.Value(loggerContextKey).(*slog.Logger); ok && log != nil {
		return log
	}
	return slog.Default()
}

// Logger returns a middleware that injects the logger into the request context and logs each response with metadata.
// Response status and bytes written are stored in context (updated by a response writer wrapper) for use in logging and downstream code.
func Logger(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			status := http.StatusOK
			written := int64(0)
			ctx := r.Context()
			ctx = context.WithValue(ctx, loggerContextKey, log)
			ctx = context.WithValue(ctx, responseStatusCodeContextKey, &status)
			ctx = context.WithValue(ctx, responseBytesWrittenContextKey, &written)
			r = r.WithContext(ctx)

			reqID := r.Header.Get("X-Request-Id")
			if reqID == "" {
				reqID = "-"
			}
			attrs := []slog.Attr{
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("request_id", reqID),
			}

			ww := &responseWriter{ResponseWriter: w, status: &status, written: &written}
			next.ServeHTTP(ww, r)

			attrs = append(attrs,
				slog.Int("status", ResponseStatusFromContext(r.Context())),
				slog.Int64("bytes", ResponseBytesWrittenFromContext(r.Context())),
				slog.Duration("duration_ms", time.Since(start)),
			)
			log.LogAttrs(r.Context(), slog.LevelInfo, "request", attrs...)
		})
	}
}

func (w *responseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	*w.written += int64(n)
	return n, err
}

// ContentTypeJSON sets Content-Type: application/json on all responses.
func ResponseContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// RequireAcceptJSON rejects requests whose Accept header does not allow application/json.
// Accepts: application/json, */*, or no Accept header.
func RequireAcceptJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept")
		if accept == "" {
			next.ServeHTTP(w, r)
			return
		}
		// Allow application/json or */*
		if strings.Contains(accept, "application/json") || strings.Contains(accept, "*/*") {
			next.ServeHTTP(w, r)
			return
		}
		httpErrors.SetHTTPError(r.Context(), httpErrors.HTTPError(http.StatusNotAcceptable, errors.New("Accept must be application/json")))
	})
}

// RequireContentTypeJSONForBody rejects POST, PUT, PATCH requests that do not have Content-Type: application/json.
// Returns 415 Unsupported Media Type. GET and other methods without a body are not checked.
func RequireContentTypeJSONForBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost && r.Method != http.MethodPut && r.Method != http.MethodPatch {
			next.ServeHTTP(w, r)
			return
		}
		ct := r.Header.Get("Content-Type")
		if strings.TrimSpace(strings.Split(ct, ";")[0]) == "application/json" {
			next.ServeHTTP(w, r)
			return
		}
		httpErrors.SetHTTPError(r.Context(), httpErrors.HTTPError(http.StatusUnsupportedMediaType, errors.New("Content-Type must be application/json")))
	})
}
