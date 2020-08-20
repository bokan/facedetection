package download

import (
	"context"
	"io"
)

type Downloader interface {
	Download(ctx context.Context, url string) (io.ReadCloser, error)
}
