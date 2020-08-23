package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_run(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "happy path",
			args:    []string{"facedetection", "-c", "../../pkg/facedetect/pigofacedetect/cascades", "-p", "0"},
			wantErr: false,
		},
		{
			name:    "make flag parse fail",
			args:    []string{"facedetection", "-x"},
			wantErr: true,
		},
		{
			name:    "make pigo cascade load fail",
			args:    []string{"facedetection", "-c", "/wrongdir", "-p", "0"},
			wantErr: true,
		},
		{
			name:    "make Serve() fail",
			args:    []string{"facedetection", "-c", "../../pkg/facedetect/pigofacedetect/cascades", "-p", "80000"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*20)
			_ = cancel
			err := run(ctx, tt.args, output)
			if (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// TODO: Write better test
func Test_requestLogger(t *testing.T) {
	output := &bytes.Buffer{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte("{}"))
	})
	rl := requestLogger(initLogger(output))

	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/foo", nil)
	rl(handler).ServeHTTP(rec, r)
}

func Test_locateCascades(t *testing.T) {
	if locateCascades() == "" {
		t.Error("locateCascades() should return a path")
	}
}
