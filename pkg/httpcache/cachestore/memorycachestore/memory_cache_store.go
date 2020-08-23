package memorycachestore

import (
	"sync"

	"github.com/bokan/facedetection/pkg/httpcache/cachestore"
)

// MemoryCacheStore is in-memory store for HTTPCache.
//
// It's implemented with a sync.Map. This implementation does not expire cache entries
// and cache size is not limited.
type MemoryCacheStore struct {
	store sync.Map
}

// NewMemoryCacheStore instantiates new MemoryCacheStore.
func NewMemoryCacheStore() *MemoryCacheStore {
	return &MemoryCacheStore{}
}

// Save saves a cache entry to store.
func (m *MemoryCacheStore) Save(key string, response *cachestore.Response) error {
	m.store.Store(key, *response)
	return nil
}

// Load retrieves a cache entry from the store.
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
