package pigofacedetect

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	"github.com/bokan/facedetection/pkg/facedetect"
)

var (
	update = flag.Bool("update", false, "update the golden files of this test")
)

func TestPigoFaceDetect_DetectFaces(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		wantErr  bool
	}{
		{
			name:     "People001",
			fileName: "people001.jpg",
			wantErr:  false,
		},
		{
			name:     "People001 Truncated to 50k",
			fileName: "people001trunc.jpg",
			wantErr:  true,
		},
		{
			name:     "People002",
			fileName: "people002.jpg",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfd := NewPigoFaceDetector()
			if err := pfd.LoadCascades("cascades"); err != nil {
				t.Fatal(err)
			}
			f, err := os.Open(path.Join("testdata", tt.fileName))
			if err != nil {
				t.Error(err)
			}

			faces, err := pfd.DetectFaces(context.Background(), f)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectFaces() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, err := json.Marshal(&faces)
			if err != nil {
				t.Fatalf("unable to marshal response: %v", err)
				return
			}

			gfPath := path.Join("testdata", tt.fileName+".json")

			if *update == true {
				if err := ioutil.WriteFile(gfPath, got, 0644); err != nil {
					t.Fatalf("unable to write golden file: %v", err)
					return
				}
				t.Skip("golden file updated")
				return
			}

			want, err := ioutil.ReadFile(gfPath)
			if err != nil {
				t.Fatalf("unable to read golden file: %v", err)
				return
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("DetectFaces() got = %v, want %v", string(got), string(want))
			}
		})
	}
}

func TestPigoFaceDetect_DetectFacesWithInvalidFormat(t *testing.T) {
	pfd := NewPigoFaceDetector()
	if err := pfd.LoadCascades("cascades"); err != nil {
		t.Fatal(err)
	}
	_, err := pfd.DetectFaces(context.Background(), strings.NewReader("foo"))
	if err != facedetect.ErrImageError {
		t.Error("detect faces should return facedetect.ErrImageError when unsupported image format is supplied")
	}

}
