package fakedownloader

import (
	"context"
	"io"
)

type FakeDownloader struct {
	rc  io.ReadCloser
	err error
}

func NewFakeDownloader(rc io.ReadCloser, err error) *FakeDownloader {
	return &FakeDownloader{rc: rc, err: err}
}

func (f FakeDownloader) Download(ctx context.Context, url string) (io.ReadCloser, error) {
	return f.rc, f.err
}
