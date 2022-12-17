package zaphttp

import (
	"net/http"

	"github.com/felixge/httpsnoop"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Handler(routerName string, l *zap.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m := httpsnoop.CaptureMetrics(h, w, r)

			lvl := zapcore.InfoLevel
			if m.Code >= 500 {
				lvl = zapcore.ErrorLevel
			}

			l.Log(lvl, "request",
				zap.String("router", routerName),
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.Int("status", m.Code),
				zap.Duration("latency", m.Duration),
				zap.Int64("bytes", m.Written),
			)
		})
	}
}
