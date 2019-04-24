package lib

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
)

func identify(root string, rec *contrib.LBPHFaceRecognizer) error {
	info, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}

	for _, i := range info {
		if i.IsDir() || i.Name() == ".DS_Store" || strings.HasSuffix(i.Name(), ".iso") {
			continue
		}

		reader, err := ioutil.ReadFile(root + i.Name())
		if err != nil {
			return err
		}

		mat, err := gocv.IMDecode(reader, gocv.IMReadGrayScale)
		if err != nil {
			return err
		}

		conf := rec.PredictExtendedResponse(mat).Confidence
		if conf <= 30.0 {
			fmt.Printf("Face Found at: %s, with error rate: %f\n", root+i.Name(), conf)
		}

		mat.Close()
	}

	return nil
}
