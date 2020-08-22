package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bokan/facedetection/pkg/api"
	"github.com/bokan/facedetection/pkg/download/httpdownloader"
	"github.com/bokan/facedetection/pkg/facedetect/pigofacedetect"
	"github.com/bokan/facedetection/pkg/httpcache"
	"github.com/bokan/facedetection/pkg/httpcache/cachestore/memorycachestore"
	"github.com/bokan/facedetection/pkg/requestlog"
	"go.uber.org/zap"
)

const (
	ExitCodeError = 1
	MaxFileSize   = 1 << 21 // 2 MiB
)

func run(ctx context.Context, args []string) error {
	log, err := initLogger()
	if err != nil {
		return err
	}

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	var (
		port         = flags.Int("p", 8000, "configure listen port")
		cascadesPath = flags.String("c", "pkg/facedetect/pigofacedetect/cascades", "configure cascades path")
	)
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	d := httpdownloader.NewHTTPDownloader(http.DefaultClient, time.Second*5, MaxFileSize)
	fd := pigofacedetect.NewPigoFaceDetector()
	if err := fd.LoadCascades(*cascadesPath); err != nil {
		log.Fatalw("PigoFaceDetector was unable to load cascades, provide cascade dir with -c flag", "dir", *cascadesPath)
		return err
	}
	a := api.NewAPI(fmt.Sprintf(":%d", *port), d, fd)

	cache := httpcache.NewHTTPCache(memorycachestore.NewMemoryCacheStore()).Middleware()
	rl := requestLogger(log)

	log.Infow("Starting service", "port", *port)
	if err := a.Serve(ctx, rl(cache(a.Routes()))); err != nil {
		if err == http.ErrServerClosed {
			log.Warn("Context ended, server stopped.")
			return nil
		}
		log.Errorw("Serve() error", "err", err)
		return err
	}
	return nil
}

func initLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	sugar := logger.Sugar()
	return sugar, nil
}

func requestLogger(log *zap.SugaredLogger) func(handler http.Handler) http.Handler {
	rl := requestlog.NewRequestLogger(func(kv map[string]interface{}) {
		var args []interface{}
		for k, v := range kv {
			args = append(args, k)
			args = append(args, v)
		}
		log.Debugw("Request", args...)
	}).Middleware()
	return rl
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		select {
		case <-sigc:
			cancel()
		case <-ctx.Done():
		}
	}()

	if err := run(ctx, os.Args); err != nil {
		os.Exit(ExitCodeError)
	}

	select {
	case <-sigc:
	case <-time.After(time.Second):
	}
}
