package requestlog

import (
	"net/http"
	"time"

	"github.com/bokan/facedetection/pkg/responserecorder"
)

type RequestLoggerFunc func(kv map[string]interface{})

type RequestLogger struct {
	lf RequestLoggerFunc
}

func NewRequestLogger(lf RequestLoggerFunc) *RequestLogger {
	return &RequestLogger{lf: lf}
}

func (l *RequestLogger) Middleware() func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			before := time.Now()
			rr := responserecorder.NewResponseRecorder(w, r)
			handler.ServeHTTP(rr, r)
			took := time.Since(before)

			info := make(map[string]interface{})
			info["ip"] = getIP(r)
			info["method"] = r.Method
			info["url"] = r.URL.String()
			info["status"] = rr.StatusCode()
			info["took"] = took.Milliseconds()
			l.lf(info)
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
