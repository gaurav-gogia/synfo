package lib

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/h2non/filetype"
	"gocv.io/x/gocv"
)

// Verify function is the entry point for face recognition module
// poidir -> Path to POI images
// extdir -> Path to extracted images from evidence
func Verify(poidir, extdir string) error {
	fmt.Println("Preparing Network ....")
	net := gocv.ReadNet(MODEL, CONFIG)
	defer net.Close()
	if err := prepareNet(&net); err != nil {
		return err
	}

	fmt.Println("Grabbing Faces ....")
	if err := detect(poidir, TRAINPOI, &net); err != nil {
		return err
	}

	if err := detect(extdir, TESTPOI, &net); err != nil {
		return err
	}

	fmt.Println("Training ....")
	rec, err := train(TRAINPOI)
	if err != nil {
		return err
	}

	fmt.Println("Identifying ....")
	err = identify(TESTPOI, rec)
	return err
}

func prepareNet(net *gocv.Net) error {
	if net.Empty() {
		return errors.New("error reading model")
	}

	backend := gocv.NetBackendDefault
	target := gocv.NetTargetCPU

	net.SetPreferableBackend(gocv.NetBackendType(backend))
	net.SetPreferableTarget(gocv.NetTargetType(target))

	return nil
}

func confirm(path string, finfo os.FileInfo) bool {
	if finfo.IsDir() || finfo.Name() == ".DS_Store" {
		return false
	}

	buff, _ := ioutil.ReadFile(path + finfo.Name())
	if filetype.IsImage(buff) {
		return true
	}

	return false
}
