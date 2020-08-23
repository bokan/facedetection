package httpdownloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bokan/facedetection/pkg/download"
)

// HTTPDownloader downloads a file from an HTTP server.
type HTTPDownloader struct {
	client        *http.Client
	clientTimeOut time.Duration
	maxFileSize   int64
}

// NewHTTPDownloader instantiates a new HTTPDownloader.
//
// Requests will be time bound by clientTimeOut.
// Files bigger than maxFileSize will be rejected.
func NewHTTPDownloader(client *http.Client, clientTimeOut time.Duration, maxFileSize int64) *HTTPDownloader {
	return &HTTPDownloader{clientTimeOut: clientTimeOut, maxFileSize: maxFileSize, client: client}
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
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, download.ErrNon200StatusCode
	}
	if resp.ContentLength > d.maxFileSize {
		return nil, download.ErrFileIsTooBig
	}

	return resp.Body, nil
}
