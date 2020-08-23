package httpcache

import (
	"fmt"
	"net/http"

	"github.com/bokan/facedetection/pkg/httpcache/cachestore"
	"github.com/bokan/facedetection/pkg/responserecorder"
)

// HTTPCache caches the successful HTTP responses.
type HTTPCache struct {
	store cachestore.CacheStore
}

// NewHTTPCache instantiates a new HTTPCache with provided cache store.
func NewHTTPCache(store cachestore.CacheStore) *HTTPCache {
	return &HTTPCache{store: store}
}

// Middleware returns a HTTP middleware that performs caching.
func (c *HTTPCache) Middleware() func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := fmt.Sprintf("%s-%s", r.Method, r.URL.String())

			if resp, err := c.store.Load(key); err == nil {
				dst := w.Header()
				for k, vv := range resp.Header {
					for _, v := range vv {
						dst.Add(k, v)
					}
				}
				w.Header().Set("X-Cache", "HIT")
				w.WriteHeader(resp.StatusCode)
				_, _ = w.Write(resp.Body)
				return
			}

			w.Header().Set("X-Cache", "MISS")
			rr := responserecorder.NewResponseRecorder(w, r)
			handler.ServeHTTP(rr, r)
			if rr.StatusCode() == 200 {
				_ = c.store.Save(key, &cachestore.Response{
					StatusCode: rr.StatusCode(),
					Header:     rr.Header(),
					Body:       rr.Body(),
				})
			}
		})
	}
}
