package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bokan/facedetection/pkg/download/fakedownloader"
	"github.com/bokan/facedetection/pkg/facedetect"
	"github.com/bokan/facedetection/pkg/facedetect/fakefacedetect"
)

func TestAPI_handleFaceDetect_WithoutImageURLParam(t *testing.T) {
	a := &API{}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	a.handleFaceDetect(rec, req)
	if rec.Result().StatusCode != 400 {
		t.Error("handler should return status code 400 when there is no image_url param")
	}
}

func TestAPI_handleFaceDetect_InvalidImageURL(t *testing.T) {
	a := &API{}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?image_url=:", nil)
	a.handleFaceDetect(rec, req)
	if rec.Result().StatusCode != 400 {
		t.Error("handler should return status code 400 when image_url contains invalid url")
	}
}

func TestAPI_handleFaceDetect_InvalidImageURLScheme(t *testing.T) {
	a := &API{}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?image_url=foo://", nil)
	a.handleFaceDetect(rec, req)
	if rec.Result().StatusCode != 400 {
		t.Error("handler should return status code 400 when image_url contains invalid url with invalid scheme")
	}
}

func TestAPI_handleFaceDetect_DownloaderError(t *testing.T) {
	a := &API{
		d: fakedownloader.NewFakeDownloader(nil, fmt.Errorf("fake error")),
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?image_url=http://localhost/", nil)
	a.handleFaceDetect(rec, req)
	if rec.Result().StatusCode != 400 {
		t.Error("handler should return status code 400 when Downloader returns an error")
	}
}

func TestAPI_handleFaceDetect_FaceDetectorErrorUnsupportedImageFormat(t *testing.T) {
	a := &API{
		d:  fakedownloader.NewFakeDownloader(ioutil.NopCloser(strings.NewReader("")), nil),
		fd: fakefacedetect.NewFakeFaceDetect(nil, facedetect.ErrUnsupportedImageFormat),
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?image_url=http://localhost/", nil)
	a.handleFaceDetect(rec, req)
	if rec.Result().StatusCode != 400 {
		t.Error("handler should return status code 400 when face detector is unable to decode image")
	}
}

func TestAPI_handleFaceDetect_FaceDetectorError(t *testing.T) {
	a := &API{
		d:  fakedownloader.NewFakeDownloader(ioutil.NopCloser(strings.NewReader("")), nil),
		fd: fakefacedetect.NewFakeFaceDetect(nil, fmt.Errorf("fake error")),
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?image_url=http://localhost/", nil)
	a.handleFaceDetect(rec, req)
	if rec.Result().StatusCode != 500 {
		t.Error("handler should return status code 500 when face detector returns an error")
	}
}

func TestAPI_handleFaceDetect_Success(t *testing.T) {
	faces := []facedetect.Face{
		{
			Bounds: &facedetect.Bounds{
				X:      10,
				Y:      20,
				Height: 100,
				Width:  100,
			},
			Mouth: &facedetect.Point{
				X: 60,
				Y: 100,
			},
			RightEye: &facedetect.Point{
				X: 30,
				Y: 50,
			},
			LeftEye: &facedetect.Point{
				X: 90,
				Y: 50,
			},
		},
	}
	a := &API{
		d:  fakedownloader.NewFakeDownloader(ioutil.NopCloser(strings.NewReader("")), nil),
		fd: fakefacedetect.NewFakeFaceDetect(faces, nil),
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?image_url=http://localhost/", nil)
	a.handleFaceDetect(rec, req)
	if rec.Result().StatusCode != 200 {
		t.Error("handler should return status code 200 on success")
	}

	if rec.Result().Header.Get("Content-Type") != "application/json" {
		t.Error("handler should return content type application/json header on success")
	}

	got := rec.Body.String()
	expected := "{\"Faces\":[{\"bounds\":{\"x\":10,\"y\":20,\"height\":100,\"width\":100},\"mouth\":{\"x\":60,\"y\":100},\"right_eye\":{\"x\":30,\"y\":50},\"left_eye\":{\"x\":90,\"y\":50}}]}"
	if got != expected {
		t.Error("handler should return expected json payload")
	}

}
