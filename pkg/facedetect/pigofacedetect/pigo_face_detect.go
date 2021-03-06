package pigofacedetect

import (
	"context"
	"image"
	_ "image/jpeg" // Add JPEG support.
	_ "image/png"  // Add PNG support.
	"io"
	"io/ioutil"
	"path"

	"github.com/bokan/facedetection/pkg/facedetect"
	pigo "github.com/esimov/pigo/core"
	"github.com/fogleman/gg"
)

const (
	minSize                  = 20
	maxSize                  = 1000
	shiftFactor              = 0.1
	scaleFactor              = 1.1
	angle                    = 0.0
	ioUThreshold             = 0.2
	perturbs                 = 63
	featuresQualityThreshold = 5.0
)

var (
	mouthCascade = []string{"lp93", "lp84", "lp82", "lp81"}
)

// PigoFaceDetector is face detector implementation based on pigo.
type PigoFaceDetector struct {
	classifier *pigo.Pigo
	plc        *pigo.PuplocCascade
	flpcs      map[string][]*pigo.FlpCascade
}

// NewPigoFaceDetector creates a new instance of PigoFaceDetector, a pigo based
// implementation of FaceDetector.
func NewPigoFaceDetector() *PigoFaceDetector {
	return &PigoFaceDetector{}
}

// LoadCascades loads binary cascade files required by pigo.
func (pfd *PigoFaceDetector) LoadCascades(cascadeDir string) error {
	cascadeFile, err := ioutil.ReadFile(path.Join(cascadeDir, "facefinder"))
	if err != nil {
		return err
	}

	p := pigo.NewPigo()
	// Unpack the binary file. This will return the number of cascade trees,
	// the tree depth, the threshold and the prediction from tree's leaf nodes.
	pfd.classifier, err = p.Unpack(cascadeFile)
	if err != nil {
		return err
	}

	pl := pigo.NewPuplocCascade()

	cascade, err := ioutil.ReadFile(path.Join(cascadeDir, "puploc"))
	if err != nil {
		return err
	}
	pfd.plc, err = pl.UnpackCascade(cascade)
	if err != nil {
		return err
	}

	pfd.flpcs, err = pl.ReadCascadeDir(path.Join(cascadeDir, "lps"))
	if err != nil {
		return err
	}
	return nil
}

// DetectFaces analyzes image provided by img parameter and
// returns slice of detected faces with facial features.
func (pfd PigoFaceDetector) DetectFaces(ctx context.Context, img io.Reader) ([]facedetect.Face, error) {
	_ = ctx // Reserved for possible future use.

	var dc *gg.Context
	var imgParams *pigo.ImageParams

	src, err := pigo.DecodeImage(img)
	if err != nil {
		if err == image.ErrFormat {
			return nil, facedetect.ErrImageError
		}
		return nil, err
	}

	pixels := pigo.RgbToGrayscale(src)
	cols, rows := src.Bounds().Max.X, src.Bounds().Max.Y

	dc = gg.NewContext(cols, rows)
	dc.DrawImage(src, 0, 0)

	imgParams = &pigo.ImageParams{
		Pixels: pixels,
		Rows:   rows,
		Cols:   cols,
		Dim:    cols,
	}

	cParams := pigo.CascadeParams{
		MinSize:     minSize,
		MaxSize:     maxSize,
		ShiftFactor: shiftFactor,
		ScaleFactor: scaleFactor,
		ImageParams: *imgParams,
	}

	// Run the classifier over the obtained leaf nodes and return the detection results.
	// The result contains quadruplets representing the row, column, scale and detection score.
	faces := pfd.classifier.RunCascade(cParams, angle)

	// Calculate the intersection over union (IoU) of two clusters.
	faces = pfd.classifier.ClusterDetections(faces, ioUThreshold)

	var outFaces []facedetect.Face

	for _, face := range faces {
		if face.Q < 20 {
			continue
		}
		outFace := facedetect.Face{
			Bounds: &facedetect.Bounds{
				X:      face.Col - face.Scale/2,
				Y:      face.Row - face.Scale/2,
				Height: face.Scale,
				Width:  face.Scale,
			},
		}
		if face.Q > featuresQualityThreshold && face.Scale > 50 {
			leftEye, rightEye := detectEyes(pfd.plc, face, imgParams)
			if leftEye != nil {
				outFace.LeftEye = &facedetect.Point{
					X: leftEye.Col,
					Y: leftEye.Row,
				}
			}
			if rightEye != nil {
				outFace.RightEye = &facedetect.Point{
					X: rightEye.Col,
					Y: rightEye.Row,
				}
			}
			m := detectMouth(pfd.flpcs, leftEye, rightEye, imgParams)
			if m != nil {
				outFace.Mouth = m
			}
		}
		outFaces = append(outFaces, outFace)
	}

	return outFaces, nil
}

func detectEyes(plc *pigo.PuplocCascade, face pigo.Detection, imgParams *pigo.ImageParams) (*pigo.Puploc, *pigo.Puploc) {
	// Left eye
	puploc := &pigo.Puploc{
		Row:      face.Row - int(0.075*float32(face.Scale)),
		Col:      face.Col - int(0.175*float32(face.Scale)),
		Scale:    float32(face.Scale) * 0.25,
		Perturbs: perturbs,
	}
	leftEye := plc.RunDetector(*puploc, *imgParams, angle, false)

	// Right eye
	puploc = &pigo.Puploc{
		Row:      face.Row - int(0.075*float32(face.Scale)),
		Col:      face.Col + int(0.185*float32(face.Scale)),
		Scale:    float32(face.Scale) * 0.25,
		Perturbs: perturbs,
	}
	rightEye := plc.RunDetector(*puploc, *imgParams, angle, false)

	return leftEye, rightEye
}

func detectMouth(flpcs map[string][]*pigo.FlpCascade, leftEye *pigo.Puploc, rightEye *pigo.Puploc, imgParams *pigo.ImageParams) *facedetect.Point {
	mouthMinX := maxSize
	mouthMaxX := 0
	mouthMinY := maxSize
	mouthMaxY := 0

	for _, mouth := range mouthCascade {
		for _, flpc := range flpcs[mouth] {
			flp := flpc.FindLandmarkPoints(leftEye, rightEye, *imgParams, perturbs, false)
			if flp.Row > 0 && flp.Col > 0 {
				if flp.Col > mouthMaxX {
					mouthMaxX = flp.Col
				}
				if flp.Col < mouthMinX {
					mouthMinX = flp.Col
				}
				if flp.Row > mouthMaxY {
					mouthMaxY = flp.Row
				}
				if flp.Row < mouthMinY {
					mouthMinY = flp.Row
				}
			}
		}
	}

	if mouthMinX == maxSize {
		return nil
	}
	mx := mouthMinX + (mouthMaxX-mouthMinX)/2
	my := mouthMinY + (mouthMaxY-mouthMinY)/2
	return &facedetect.Point{
		X: mx,
		Y: my,
	}
}
