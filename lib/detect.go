package lib

import (
	"image"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"

	"gocv.io/x/gocv"
)

func detect(root, savedir string, net *gocv.Net) error {
	info, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}

	for _, file := range info {
		if ok := confirm(root, file); !ok {
			continue
		}

		img, err := readfile(root, file)
		if err != nil {
			return err
		}

		blob := gocv.BlobFromImage(img, 1.0, image.Pt(300, 300), gocv.NewScalar(1.0, 1.0, 1.0, 1.0), false, false)
		defer blob.Close()
		net.SetInput(blob, "data")
		detBlob := net.Forward("detection_out")
		defer detBlob.Close()
		detections := gocv.GetBlobChannel(detBlob, 0, 0)
		defer detections.Close()

		name := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())) + ".jpeg"

		save(savedir+name, &detections, &img)
		img.Close()
	}

	return nil
}

func readfile(root string, file os.FileInfo) (gocv.Mat, error) {
	var img gocv.Mat

	reader, err := ioutil.ReadFile(root + file.Name())
	if err != nil {
		return img, err
	}

	img, err = gocv.IMDecode(reader, gocv.IMReadAnyColor)
	if err != nil {
		return img, err
	}

	return img, nil
}

func save(name string, detections, img *gocv.Mat) {
	w := float64(img.Cols())
	h := float64(img.Rows())

	for r := 0; r < detections.Rows(); r++ {
		confidence := detections.GetFloatAt(r, 2)
		if confidence < 0.25 {
			continue
		}

		left := float64(detections.GetFloatAt(r, 3)) * w
		top := float64(detections.GetFloatAt(r, 4)) * h
		right := float64(detections.GetFloatAt(r, 5)) * w
		bottom := float64(detections.GetFloatAt(r, 6)) * h

		left = math.Min(math.Max(0.0, left), w-1.0)
		right = math.Min(math.Max(0.0, right), w-1.0)
		bottom = math.Min(math.Max(0.0, bottom), h-1.0)
		top = math.Min(math.Max(0.0, top), h-1.0)

		rect := image.Rect(int(left), int(top), int(right), int(bottom))
		mat := img.Region(rect)

		gocv.IMWrite(name, mat)

		mat.Close()
	}
}
