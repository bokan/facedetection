package memorycachestore

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/bokan/facedetection/pkg/httpcache/cachestore"
)

func TestMemoryCacheStore(t *testing.T) {

	mcs := NewMemoryCacheStore()
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	resp := cachestore.Response{
		StatusCode: 200,
		Header:     header,
		Body:       []byte("bar"),
	}
	key := "foo"
	if err := mcs.Save(key, &resp); err != nil {
		t.Errorf("Save() should not return an error, got: %v", err)
		return
	}

	got, err := mcs.Load(key)
	if err != nil {
		t.Errorf("Save() should not return an error, got: %v", err)
	}

	if !reflect.DeepEqual(resp, *got) {
		t.Errorf("Response from Load() does not match with one we used with Save(), want = %v, got = %v", resp, *got)
	}
}

func TestMemoryCacheStore_Load_CacheMiss(t *testing.T) {
	mcs := NewMemoryCacheStore()
	key := "foo"

	_, err := mcs.Load(key)
	if err != cachestore.ErrCacheMiss {
		t.Errorf("Load() on empty store should return ErrCacheMiss")
	}
}

func TestMemoryCacheStore_Load_InvalidResponse(t *testing.T) {
	mcs := NewMemoryCacheStore()
	key := "foo"

	mcs.store.Store(key, "bar")

	_, err := mcs.Load(key)
	if err != cachestore.ErrInvalidCacheResponse {
		t.Errorf("Load() of invalid cache entry should return ErrInvalidCacheResponse")
	}
}
