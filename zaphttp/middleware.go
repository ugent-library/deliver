package zaphttp

import (
	"context"
	"net/http"

	"github.com/felixge/httpsnoop"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var loggerKey = contextKey("logger")

func Logger(c context.Context) *zap.Logger {
	if l := c.Value(loggerKey); l != nil {
		return l.(*zap.Logger)
	}
	return nil
}

func SetLogger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := logger
			if requestID := r.Header.Get("X-Request-ID"); requestID != "" {
				l = l.With(zap.String("requestID", requestID))
			}
			c := context.WithValue(r.Context(), loggerKey, l)
			next.ServeHTTP(w, r.WithContext(c))
		})
	}
}

func LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := Logger(r.Context())

		if l == nil {
			next.ServeHTTP(w, r)
			return
		}

		m := httpsnoop.CaptureMetrics(next, w, r)

		lvl := zapcore.InfoLevel
		if m.Code >= 500 {
			lvl = zapcore.ErrorLevel
		}

		l.Log(lvl, "request",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.Int("status", m.Code),
			zap.Duration("latency", m.Duration),
			zap.Int64("bytes", m.Written),
		)
	})
}
