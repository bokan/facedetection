package requestlog

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRequestLogger_Middleware(t *testing.T) {
	want := make(map[string]interface{})
	want["ip"] = "127.0.0.1"
	want["method"] = http.MethodGet
	want["url"] = "/foo"
	want["status"] = 200

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(want["status"].(int))
		_, _ = w.Write([]byte("{}"))
	})
	var got map[string]interface{}
	rl := NewRequestLogger(func(kv map[string]interface{}) {
		got = kv
	})
	m := rl.Middleware()(handler)

	rec := httptest.NewRecorder()
	r := httptest.NewRequest(want["method"].(string), want["url"].(string), nil)
	r.RemoteAddr = want["ip"].(string)

	m.ServeHTTP(rec, r)

	for k := range want {
		if !reflect.DeepEqual(want[k], got[k]) {
			t.Errorf("expected %s = !%v!, got = !%v!", k, want[k], got[k])
		}
	}
}

func Test_getIP_NoHeader(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/foo", nil)
	want := "127.0.0.1"
	r.RemoteAddr = want
	got := getIP(r)
	if got != want {
		t.Errorf("getIp() returned wrong ip, want = %v, got = %v", want, got)
	}
}
func Test_getIP_XForwardedFor(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/foo", nil)
	r.RemoteAddr = "127.0.0.1"
	want := "10.0.0.1"
	r.Header.Set("X-Forwarded-For", want)
	got := getIP(r)
	if got != want {
		t.Errorf("getIp() returned wrong ip, want = %v, got = %v", want, got)
	}
}
