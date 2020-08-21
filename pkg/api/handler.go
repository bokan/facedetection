package api

import (
	"encoding/json"
	"image"
	"net/http"
	"net/url"

	"github.com/bokan/stream/pkg/facedetect"
)

type Error struct {
	Message string `json:"message"`
}

// Faces structure is response sent to client. It encapsulates response from face detector.
type Faces struct {
	Faces []facedetect.Face `json:"Faces"`
}

func (a *API) handleFaceDetect(w http.ResponseWriter, r *http.Request) {
	imageURL, ok := r.URL.Query()["image_url"]
	if !ok || len(imageURL) == 0 {
		http.Error(w, "image_url query parameter missing", 400)
		return
	}

	u, err := url.Parse(imageURL[0])
	if err != nil {
		http.Error(w, "image_url is not a valid url", 400)
		return
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		http.Error(w, "image_url scheme must be http or https", 400)
		return
	}

	body, err := a.d.Download(r.Context(), imageURL[0])
	if err != nil {
		http.Error(w, "image download failed", 400)
		return
	}
	defer func() {
		_ = body.Close()
	}()

	detections, err := a.fd.DetectFaces(r.Context(), body)
	if err != nil {
		if err == image.ErrFormat {
			http.Error(w, "unsupported image format", 400)
			return
		}
		http.Error(w, "an internal error happened the during face detection", 500)
		return
	}

	response := Faces{Faces: detections}
	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "an internal error happened", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, _ = w.Write(js)
}
