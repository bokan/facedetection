package download

import (
	"context"
	"fmt"
	"io"
)

var (
	// ErrNon200StatusCode is returned by Downloader.Download calls when server returns an error code different than 200.
	ErrNon200StatusCode = fmt.Errorf("server did not return status code 200")

	// ErrFileIsTooBig is returned by Downloader.Download calls when requested file size is too big.
	ErrFileIsTooBig = fmt.Errorf("file is too big")
)

// Downloader initiates a download of a resource specified by url parameter.
type Downloader interface {
	Download(ctx context.Context, url string) (io.ReadCloser, error)
}
