package requestlog

import (
	"net/http"
	"time"

	"github.com/bokan/facedetection/pkg/responserecorder"
	"go.uber.org/zap"
)

type RequestLogger struct {
	log *zap.SugaredLogger
}

func NewRequestLogger(log *zap.SugaredLogger) *RequestLogger {
	return &RequestLogger{log: log}
}

func (l *RequestLogger) Middleware() func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			before := time.Now()
			rr := responserecorder.NewResponseRecorder(w, r)
			handler.ServeHTTP(rr, r)
			took := time.Since(before)
			l.log.Debugw("Request", "ip", getIP(r), "method", r.Method, "url", r.URL, "status", rr.StatusCode(), "took", took.Milliseconds())
		})
	}
}

func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
