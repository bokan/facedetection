package integration

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/bokan/facedetection/pkg/api"
	"github.com/bokan/facedetection/pkg/download/httpdownloader"
	"github.com/bokan/facedetection/pkg/facedetect/pigofacedetect"
)

var (
	update = flag.Bool("update", false, "update the golden files of this test")
)

func TestAPI(t *testing.T) {
	fs := httptest.NewServer(http.FileServer(http.Dir("testdata")))

	d := httpdownloader.NewHTTPDownloader(http.DefaultClient, time.Second, 1<<20)
	fd := pigofacedetect.NewPigoFaceDetector()
	if err := fd.LoadCascades("../pkg/facedetect/pigofacedetect/cascades"); err != nil {
		t.Fatal(err)
	}
	a := api.NewAPI("", d, fd)

	apis := httptest.NewServer(a.Routes())
	fn := "people001.jpg"
	imgURL := fmt.Sprintf("%s/%s", fs.URL, fn)
	url := fmt.Sprintf("%s/v1/face-detect?image_url=%s", apis.URL, imgURL)
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Error(err)
	}
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		t.Error(err)
	}

	got := make([]byte, 4096)
	if _, err := resp.Body.Read(got); err != nil {
		if err != io.EOF {
			t.Error(err)
		}
	}

	gf := path.Join("testdata", fn+".json")

	if *update == true {
		if err := ioutil.WriteFile(gf, got, 0644); err != nil {
			t.Fatalf("unable to write golden file: %v", err)
			return
		}
		t.Skip("golden file updated")
		return
	}

	want, err := ioutil.ReadFile(gf)
	if err != nil {
		t.Fatalf("unable to read golden file: %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", string(got), string(want))
	}
}
