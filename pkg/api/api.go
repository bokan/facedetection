package api

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/bokan/stream/pkg/download"
	"github.com/bokan/stream/pkg/facedetect"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type API struct {
	addr string
	d    download.Downloader
	fd   facedetect.FaceDetector
	srv  http.Server
}

// NewAPI creates a REST API responsible for serving face detection requests.
// Call Serve afterwards and pass parent context and return value of Routes
// as parameters,
func NewAPI(addr string, d download.Downloader, fd facedetect.FaceDetector) *API {
	return &API{addr: addr, d: d, fd: fd}
}

// Routes returns a http.Handler with routes configured.
func (a *API) Routes() http.Handler {
	r := mux.NewRouter()
	r.Methods(http.MethodGet).Path("/v1/face-detect").HandlerFunc(a.handleFaceDetect)

	allowAllOrigins := handlers.AllowedOriginValidator(func(origin string) bool {
		return true // Allow all origins
	})
	headersOk := handlers.AllowedHeaders([]string{})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "OPTIONS"})

	return handlers.CORS(headersOk, allowAllOrigins, methodsOk)(r)
}

// Serve starts a HTTP server and serves provided handler. To invoke face detection
// endpoint, perform a GET request on /v1/face-detect?={image_url}
func (a *API) Serve(ctx context.Context, handler http.Handler) error {
	a.srv = http.Server{
		Addr:              a.addr,
		ReadTimeout:       time.Second * 2,
		ReadHeaderTimeout: time.Second * 2,
		WriteTimeout:      time.Second * 5,
		IdleTimeout:       time.Second * 5,
		MaxHeaderBytes:    1024,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
		Handler: handler,
	}
	go func() {
		<-ctx.Done()
		_ = a.srv.Shutdown(ctx)
	}()
	if err := a.srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
