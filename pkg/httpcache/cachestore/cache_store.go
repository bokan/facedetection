package cachestore

import (
	"errors"
	"net/http"
)

var (
	ErrCacheMiss            = errors.New("cache miss")
	ErrInvalidCacheResponse = errors.New("invalid cache response")
)

type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

type CacheStore interface {
	Save(key string, response *Response) error
	Load(key string) (*Response, error)
}
