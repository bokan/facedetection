package download

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDownloader_Download1024ByteFile(t *testing.T) {
	d := NewDownloader(time.Second*5, 4*1024*1024)
	ctx := context.Background()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
		buf := make([]byte, 1024)
		_, _ = w.Write(buf)
	}))

	got, err := d.Download(ctx, srv.URL)
	if err != nil {
		t.Errorf("Download() error = %v", err)
		return
	}
	buf := make([]byte, 1024)
	n, err := got.Read(buf)
	if err != nil {
		if err != io.EOF {
			t.Errorf("%v", err)
			return
		}
	}
	if n != 1024 {
		t.Errorf("expected 1024 bytes, got %d", n)
		return
	}
}

func TestDownloader_InvalidUrl(t *testing.T) {
	d := NewDownloader(time.Second*5, 512)
	ctx := context.Background()

	_, err := d.Download(ctx, "")
	if err == nil {
		t.Errorf("should return invalid url error")
		return
	}
}

func TestDownloader_Non200StatusCode(t *testing.T) {
	d := NewDownloader(time.Second*5, 4*1024*1024)
	ctx := context.Background()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(404)
		buf := make([]byte, 1024)
		_, _ = w.Write(buf)
	}))

	_, err := d.Download(ctx, srv.URL)
	if err == nil {
		t.Errorf("status code 404 should cause an error")
		return
	}
}

func TestDownloader_FileTooBig(t *testing.T) {
	d := NewDownloader(time.Second*5, 512)
	ctx := context.Background()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Length", "1024")
		w.WriteHeader(200)
		buf := make([]byte, 1024)
		_, _ = w.Write(buf)
	}))

	_, err := d.Download(ctx, srv.URL)
	if err == nil {
		t.Errorf("should return file too big error")
		return
	}
}
