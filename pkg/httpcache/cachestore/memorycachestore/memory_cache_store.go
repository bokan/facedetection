package memorycachestore

import (
	"sync"

	"github.com/bokan/facedetection/pkg/httpcache/cachestore"
)

type MemoryCacheStore struct {
	store sync.Map
}

func NewMemoryCacheStore() *MemoryCacheStore {
	return &MemoryCacheStore{}
}

func (m *MemoryCacheStore) Save(key string, response *cachestore.Response) error {
	m.store.Store(key, *response)
	return nil
}

func (m *MemoryCacheStore) Load(key string) (*cachestore.Response, error) {
	value, ok := m.store.Load(key)
	if !ok {
		return nil, cachestore.ErrCacheMiss
	}
	resp, valid := value.(cachestore.Response)
	if !valid {
		return nil, cachestore.ErrInvalidCacheResponse
	}

	return &resp, nil
}
