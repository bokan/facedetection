package api

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/bokan/stream/pkg/download"
	"github.com/bokan/stream/pkg/facedetect"
	"github.com/gorilla/mux"
)

type API struct {
	addr string
	d    download.Downloader
	fd   facedetect.FaceDetector
}

func NewAPI(addr string, downloader download.Downloader, fd facedetect.FaceDetector) *API {
	return &API{addr: addr, d: downloader, fd: fd}
}

func (a *API) Routes() http.Handler {
	r := mux.NewRouter()
	r.Methods(http.MethodGet).Path("/v1/face-detect").HandlerFunc(a.handleFaceDetect)
	return r
}

func (a *API) Serve(ctx context.Context) error {
	srv := http.Server{
		Addr:              a.addr,
		ReadTimeout:       time.Second * 2,
		ReadHeaderTimeout: time.Second * 2,
		WriteTimeout:      time.Second * 5,
		IdleTimeout:       time.Second * 5,
		MaxHeaderBytes:    1024,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
		Handler: a.Routes(),
	}
	if err := srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
