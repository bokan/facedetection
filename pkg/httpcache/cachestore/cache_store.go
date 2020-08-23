package cachestore

import (
	"errors"
	"net/http"
)

var (
	// ErrCacheMiss is returned by CacheStore.Load calls when there is no entry in store for a given key.
	ErrCacheMiss = errors.New("cache miss")

	// ErrInvalidCacheResponse is returned by CacheStore.Load calls in case of cache entry corruption.
	ErrInvalidCacheResponse = errors.New("invalid cache response")
)

// Response is used to store HTTP response in cache.
type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

// CacheStore acts as a storage for HTTPCache.
type CacheStore interface {
	Save(key string, response *Response) error
	Load(key string) (*Response, error)
}
