package httpdownloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTTPDownloader struct {
	clientTimeOut time.Duration
	maxFileSize   int64
	client        http.Client
}

func NewHTTPDownloader(clientTimeOut time.Duration, maxFileSize int64) *HTTPDownloader {
	return &HTTPDownloader{clientTimeOut: clientTimeOut, maxFileSize: maxFileSize, client: http.Client{}}
}

// Download initiates a time constrained HTTP GET request, validates Content-Length and returns response body io.ReadCloser.
func (d *HTTPDownloader) Download(ctx context.Context, url string) (io.ReadCloser, error) {

	// TODO: Redirects?

	ctx, cancel := context.WithTimeout(ctx, d.clientTimeOut)
	_ = cancel

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status code %d", resp.StatusCode)
	}
	if resp.ContentLength > d.maxFileSize {
		return nil, fmt.Errorf("file is too big")
	}

	return resp.Body, nil
}
