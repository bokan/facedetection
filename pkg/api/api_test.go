package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAPI_Routes_CORS(t *testing.T) {
	a := NewAPI("", nil, nil)
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

func TestAPI_Serve(t *testing.T) {
	a := NewAPI(":0", nil, nil)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	_ = cancel
	if err := a.Serve(ctx, nil); err != http.ErrServerClosed {
		t.Errorf("serve should return ErrServerClosed error when context ends, got: %v", err)
	}
}
