package fakefacedetect

import (
	"context"
	"io"

	"github.com/bokan/facedetection/pkg/facedetect"
)

type FakeFaceDetect struct {
	detections []facedetect.Face
	err        error
}

func NewFakeFaceDetect(detections []facedetect.Face, err error) *FakeFaceDetect {
	return &FakeFaceDetect{detections: detections, err: err}
}

func (f FakeFaceDetect) DetectFaces(ctx context.Context, img io.Reader) ([]facedetect.Face, error) {
	return f.detections, f.err
}
