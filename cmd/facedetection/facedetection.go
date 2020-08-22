package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/bokan/facedetection/pkg/api"
	"github.com/bokan/facedetection/pkg/download/httpdownloader"
	"github.com/bokan/facedetection/pkg/facedetect/pigofacedetect"
	"go.uber.org/zap"
)

const (
	ExitCodeError = 1
)

func run(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	var (
		port         = flags.Int("p", 8000, "-p <listen_port>")
		cascadesPath = flags.String("c", "pkg/facedetect/pigofacedetect/cascades", "-c <cascades_path>")
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
	fd := pigofacedetect.NewPigoFaceDetector()
	if err := fd.LoadCascades(*cascadesPath); err != nil {
		sugar.Fatalw("PigoFaceDetector was unable to load cascades, provide cascade dir with -c flag", "dir", *cascadesPath)
		return err
	}
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
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
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
