package facedetect

import (
	"context"
	"fmt"
	"io"
)

var (
	ErrUnsupportedImageFormat = fmt.Errorf("unsupported image format")
	ErrImageError             = fmt.Errorf("error loading image")
)

type Bounds struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Height int `json:"height"`
	Width  int `json:"width"`
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Face struct {
	Bounds   *Bounds `json:"bounds"`
	Mouth    *Point  `json:"mouth,omitempty"`
	RightEye *Point  `json:"right_eye,omitempty"`
	LeftEye  *Point  `json:"left_eye,omitempty"`
}

type FaceDetector interface {
	DetectFaces(ctx context.Context, img io.Reader) ([]Face, error)
}
