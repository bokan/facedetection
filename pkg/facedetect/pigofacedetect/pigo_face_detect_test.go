package pigofacedetect

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"testing"
)

func TestPigoFaceDetect_Foo(t *testing.T) {
	p := NewPigoFaceDetect()
	ctx := context.Background()
	f, err := os.Open("testdata/people001.jpg")
	if err != nil {
		return
	}
	faces, err := p.DetectFaces(ctx, f)
	if err != nil {
		t.Error(err)
	}

	j, err := json.Marshal(&faces)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("\nFacesJson: %v\n", string(j))

}

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
			pfd := NewPigoFaceDetect()

			f, err := os.Open(path.Join("testdata", tt.fileName))
			if err != nil {
				t.Error(err)
			}

			got, err := pfd.DetectFaces(context.Background(), f)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectFaces() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_ = got

			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("DetectFaces() got = %v, want %v", got, tt.want)
			//}
		})
	}
}
