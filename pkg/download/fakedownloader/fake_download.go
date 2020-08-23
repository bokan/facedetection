package fakedownloader

import (
	"context"
	"io"
)

// FakeDownloader is a mock Downloader.
type FakeDownloader struct {
	rc  io.ReadCloser
	err error
}

// NewFakeDownloader returns a test double downloader. Parameters passed will be returned
// by the Download function.
func NewFakeDownloader(rc io.ReadCloser, err error) *FakeDownloader {
	return &FakeDownloader{rc: rc, err: err}
}

// Download returns the parameters given to constructor.
func (f FakeDownloader) Download(ctx context.Context, url string) (io.ReadCloser, error) {
	return f.rc, f.err
}
