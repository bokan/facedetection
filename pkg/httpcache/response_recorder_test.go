package httpcache

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestResponseRecorder_WriteBeforeWriteHeaderPanic(t *testing.T) {
	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/foo", nil)
	crw := NewResponseRecorder(rec, r)
	panicked := false
	defer func() {
		if err := recover(); err != nil {
			panicked = true
		}
	}()
	_, _ = crw.Write([]byte("foo"))
	if !panicked {
		t.Error("calling Write() before WriteHeader() should panic")
		return
	}
}

func TestResponseRecorder(t *testing.T) {
	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/foo", nil)
	crw := NewResponseRecorder(rec, r)
	crw.Header().Set("Content-Length", "6")
	crw.WriteHeader(200)
	_, _ = crw.Write([]byte("foo"))
	_, _ = crw.Write([]byte("bar"))

	if crw.StatusCode() != 200 {
		t.Error("recorded status code should be 200")
	}

	if !reflect.DeepEqual(rec.Header(), crw.Header()) {
		t.Error("passthrough header does not match recorded header")
		return
	}

	if !reflect.DeepEqual(rec.Body.Bytes(), crw.buf.Bytes()) {
		t.Error("passthrough body does not match recorded body")
		return
	}

}
