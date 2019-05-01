package lib

import (
	"fmt"
	"io/ioutil"

	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
)

func identify(root string, rec *contrib.LBPHFaceRecognizer) error {
	info, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}

	for _, file := range info {
		if ok := confirm(root, file); !ok {
			continue
		}

		reader, err := ioutil.ReadFile(root + file.Name())
		if err != nil {
			return err
		}

		mat, err := gocv.IMDecode(reader, gocv.IMReadGrayScale)
		if err != nil {
			return err
		}

		conf := rec.PredictExtendedResponse(mat).Confidence
		if conf <= 30.0 {
			fmt.Printf("Face Found at: %s, with error rate: %f\n", root+file.Name(), conf)
		}

		mat.Close()
	}

	return nil
}
