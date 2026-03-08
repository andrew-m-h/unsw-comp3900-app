//go:build integration

package integration

import (
	"context"
	"log/slog"
)

const logCaptureBuffer = 128

// CapturedEntry holds one log record for assertion in tests.
type CapturedEntry struct {
	Level   slog.Level
	Message string
	Attrs   map[string]any
}

// captureHandler is an slog.Handler that sends each record on a channel for tests to drain.
type captureHandler struct {
	ch chan CapturedEntry
}

func (h *captureHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level <= slog.LevelInfo
}

func (h *captureHandler) Handle(_ context.Context, r slog.Record) error {
	attrs := make(map[string]any)
	r.Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = attrValue(a.Value)
		return true
	})
	h.ch <- CapturedEntry{Level: r.Level, Message: r.Message, Attrs: attrs}
	return nil
}

func attrValue(v slog.Value) any {
	switch v.Kind() {
	case slog.KindString:
		return v.String()
	case slog.KindInt64:
		return v.Int64()
	case slog.KindUint64:
		return v.Uint64()
	case slog.KindFloat64:
		return v.Float64()
	case slog.KindBool:
		return v.Bool()
	case slog.KindDuration:
		return v.Duration().String()
	case slog.KindTime:
		return v.Time().String()
	case slog.KindGroup:
		return v.Group()
	default:
		return v.String()
	}
}

func (h *captureHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *captureHandler) WithGroup(name string) slog.Handler {
	return h
}

// LogCapture holds the capturable handler. Each test gets its own LogCapture with its own server, so logs are isolated.
type LogCapture struct {
	*captureHandler
}

// NewLogCapture returns a LogCapture and a logger that writes to it.
func NewLogCapture() (*LogCapture, *slog.Logger) {
	h := &captureHandler{ch: make(chan CapturedEntry, logCaptureBuffer)}
	return &LogCapture{captureHandler: h}, slog.New(h)
}

// Entries drains the log channel and returns all captured entries. Call before Close (e.g. before cleanup).
func (c *LogCapture) Entries() []CapturedEntry {
	var out []CapturedEntry
	for {
		select {
		case e := <-c.ch:
			out = append(out, e)
		default:
			return out
		}
	}
}

// Close closes the log channel. Call only after the server is closed (no more requests), e.g. in cleanup.
func (c *LogCapture) Close() {
	close(c.ch)
}
