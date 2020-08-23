package facedetect

import (
	"context"
	"fmt"
	"io"
)

var (
	// ErrUnsupportedImageFormat is returned by FaceDetector.DetectFaces calls when given image format is unsupported.
	ErrUnsupportedImageFormat = fmt.Errorf("unsupported image format")

	// ErrImageError is returned by FaceDetector.DetectFaces calls when face detector is unable to decode the image.
	ErrImageError = fmt.Errorf("error loading image")
)

// Bounds contains the position and size of face boundaries rectangle.
type Bounds struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Height int `json:"height"`
	Width  int `json:"width"`
}

// Point is used to locate the facial features on an image.
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Face is a container for face boundaries and positions of facial features.
type Face struct {
	Bounds   *Bounds `json:"bounds"`
	Mouth    *Point  `json:"mouth,omitempty"`
	RightEye *Point  `json:"right_eye,omitempty"`
	LeftEye  *Point  `json:"left_eye,omitempty"`
}

// FaceDetector performs the image analysis and returns the list of detected faces with facial features.
type FaceDetector interface {
	DetectFaces(ctx context.Context, img io.Reader) ([]Face, error)
}
