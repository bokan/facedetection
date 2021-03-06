package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/bokan/facedetection/pkg/api"
	"github.com/bokan/facedetection/pkg/download/httpdownloader"
	"github.com/bokan/facedetection/pkg/facedetect/pigofacedetect"
	"github.com/bokan/facedetection/pkg/httpcache"
	"github.com/bokan/facedetection/pkg/httpcache/cachestore/memorycachestore"
	"github.com/bokan/facedetection/pkg/requestlog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	exitCodeError = 1
	maxFileSize   = 1 << 21 // 2 MiB
)

func locateCascades(goPath string) string {
	paths := strings.Split(goPath, ":")
	for _, path := range paths {
		tryPath := filepath.Join(path, "src/github.com/bokan/facedetection/pkg/facedetect/pigofacedetect/cascades")
		if _, err := os.Stat(tryPath); err == nil {
			return tryPath
		}
	}
	return ""
}

func goPath(envGoPath string) (string, error) {
	goPath := envGoPath
	if goPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		goPath = filepath.Join(home, "go")
	}
	return goPath, nil
}

func run(ctx context.Context, args []string, output io.Writer) error {
	log := initLogger(output)

	goPath, err := goPath(os.Getenv("GOPATH"))
	if err != nil {
		return err
	}
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	var (
		port         = flags.Int("p", 8000, "configure listen port")
		cascadesPath = flags.String("c", locateCascades(goPath), "configure cascades path")
	)
	flags.SetOutput(output)
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	d := httpdownloader.NewHTTPDownloader(http.DefaultClient, time.Second*5, maxFileSize)
	fd := pigofacedetect.NewPigoFaceDetector()
	if err := fd.LoadCascades(*cascadesPath); err != nil {
		log.Errorw("PigoFaceDetector was unable to load cascades, provide cascade dir with -c flag", "dir", *cascadesPath)
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

func initLogger(sync io.Writer) *zap.SugaredLogger {
	ce := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	c := zapcore.NewCore(ce, zapcore.AddSync(sync), zap.DebugLevel)
	l := zap.New(c)
	sugar := l.Sugar()
	return sugar
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

	if err := run(ctx, os.Args, os.Stdout); err != nil {
		os.Exit(exitCodeError)
	}

	select {
	case <-sigc:
	case <-time.After(time.Second):
	}
}
