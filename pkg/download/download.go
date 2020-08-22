package download

import (
	"context"
	"fmt"
	"io"
)

var (
	ErrNon200StatusCode = fmt.Errorf("server did not return status code 200")
	ErrFileIsTooBig     = fmt.Errorf("file is too big")
)

type Downloader interface {
	Download(ctx context.Context, url string) (io.ReadCloser, error)
}
