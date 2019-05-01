package lib

import (
	"io/ioutil"

	"gocv.io/x/gocv/contrib"

	"gocv.io/x/gocv"
)

func train(root string) (*contrib.LBPHFaceRecognizer, error) {
	var mats []gocv.Mat
	rec := contrib.NewLBPHFaceRecognizer()
	var ids []int

	info, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}

	for _, file := range info {
		if ok := confirm(root, file); !ok {
			continue
		}

		reader, err := ioutil.ReadFile(root + file.Name())
		if err != nil {
			return nil, err
		}

		mat, err := gocv.IMDecode(reader, gocv.IMReadGrayScale)
		if err != nil {
			return nil, err
		}

		mats = append(mats, mat)

		ids = append(ids, 0)
	}

	rec.Train(mats, ids)

	for _, mat := range mats {
		mat.Close()
	}

	return rec, nil
}
