package httpcache

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bokan/facedetection/pkg/httpcache/cachestore/memorycachestore"
)

func TestHTTPCache_Middleware(t *testing.T) {
	hc := NewHTTPCache(memorycachestore.NewMemoryCacheStore())
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte("foobar"))
	})
	m := hc.Middleware()(handler)

	recm := httptest.NewRecorder()
	rm := httptest.NewRequest(http.MethodGet, "/foo", nil)
	m.ServeHTTP(recm, rm)
	if recm.Header().Get("X-Cache") != "MISS" {
		t.Error("expected cache miss")
		return
	}

	rech := httptest.NewRecorder()
	rh := httptest.NewRequest(http.MethodGet, "/foo", nil)
	m.ServeHTTP(rech, rh)
	if rech.Header().Get("X-Cache") != "HIT" {
		t.Error("expected cache hit")
		return
	}

}
