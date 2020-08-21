package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bokan/stream/pkg/api"
	"github.com/bokan/stream/pkg/download/httpdownloader"
	"github.com/bokan/stream/pkg/facedetect/pigofacedetect"
	"go.uber.org/zap"
)

const (
	ExitCodeError = 1
)

func run(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	var (
		port = flags.Int("p", 8000, "-p <listen_port>")
	)
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	sugar := logger.Sugar()

	d := httpdownloader.NewHTTPDownloader(time.Second*5, 10*1024*1024)
	fd := pigofacedetect.NewPigoFaceDetect("pkg/facedetect/pigofacedetect/cascades")
	a := api.NewAPI(fmt.Sprintf(":%d", *port), d, fd)

	sugar.Infow("Starting service", "port", *port)
	if err := a.Serve(ctx, a.Routes()); err != nil {
		if err == http.ErrServerClosed {
			sugar.Warn("Context ended, server stopped.")
			return nil
		}
		sugar.Errorw("Serve() error", "err", err)
		return err
	}
	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	defer func() {
		signal.Stop(sigCh)
		cancel()
	}()
	go func() {
		select {
		case <-sigCh:
			cancel()
		case <-ctx.Done():
		}
	}()

	if err := run(ctx, os.Args); err != nil {
		os.Exit(ExitCodeError)
	}

}
