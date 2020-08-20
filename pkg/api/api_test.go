package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPI_Routes_CORS(t *testing.T) {
	a := &API{}
	r := a.Routes()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/face-detect", nil)
	const origin = "https://foo.bar"
	req.Header.Set("Origin", origin)
	r.ServeHTTP(rec, req)

	if rec.Header().Get("Access-Control-Allow-Origin") != origin {
		t.Error("Access-Control-Allow-Origin header missing")
	}
}
