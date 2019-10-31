package lib

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	linux = "linux"
	mac   = "darwin"
)

// Attach function mounts the iso file as a special block device
func Attach(src string) (string, string, error) {
	if runtime.GOOS == mac {
		mntpoint := genname(6)
		out, err := exec.Command("hdiutil", "attach", "-mountpoint", mntpoint, src).Output()
		if err != nil {
			return "", "", err
		}
		mntloc := strings.Fields(string(out))[0]
		copysrc := strings.Fields(string(out))[1]
		return mntloc, copysrc, nil
	} else if runtime.GOOS == linux {
		/*
			TODO:
				Mount disk image as loop device
		*/
	}

	return "", "", errors.New("unknown runtime")
}

// Detach function unmounts the attached iso file
func Detach(name string) error {
	var err error
	if runtime.GOOS == mac {
		_, err = exec.Command("hdiutil", "detach", name).Output()
		return err
	} else if runtime.GOOS == linux {
		/*
			TODO:
				Unmount the loop device
		*/
	}
	return errors.New("unknown runtime")
}

func getsize(path string) (uint64, error) {
	var size uint64

	if runtime.GOOS == mac {
		data, _ := exec.Command("diskutil", "info", path).Output()
		info := strings.Split(string(data), "\n")

		for _, str := range info {
			text := fixspace(str)
			if strings.HasPrefix(text, "Disk Size:") {
				size = getnum(text)
				break
			}
		}
		return size, nil
	} else if runtime.GOOS == linux {
		/*
			TODO:
				Get size of special block device file
		*/
	}

	return 0, errors.New("unknown runtime")
}

func genname(length int) string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base32.StdEncoding.EncodeToString(randomBytes)[:length]
}

func fixspace(data string) string {
	const (
		leadCloseExp = `^[\s\p{Zs}]+|[\s\p{Zs}]+$`
		insidersExp  = `[\s\p{Zs}]{2,}`
	)

	leadClose := regexp.MustCompile(leadCloseExp)
	insiders := regexp.MustCompile(insidersExp)
	final := leadClose.ReplaceAllString(data, "")
	return insiders.ReplaceAllString(final, " ")
}

func getnum(data string) uint64 {
	const exp = "[^0-9]+"
	num := regexp.MustCompile(exp)
	txt := num.ReplaceAllString(data, "-")
	size, _ := strconv.ParseInt(strings.Split(txt, "-")[3], 10, 64)
	return uint64(size)
}

// PyApd function runs Automated PoI Detection
// It takes following arguments:
//		poitest - directory of extracted images
//		poitrain - directory of known images with faces
//		modeltype - an option for using hog or cnn model for apd
func PyApd(poitest, poitrain, modeltype string) error {
	start := time.Now()

	fmt.Println("\n\nRunning PoI Identification  ....")
	err := exec.Command("python3", "./libpy/face.py", poitest, poitrain, modeltype).Run()

	fmt.Printf("\nPoI Identification Time: %v\n", time.Since(start))
	return err
}

// PyAwd function runs Automated Weapon Detection on images
// It takes following arguments:
//		wepimages - directory of extrcated images
func PyAwd(wepimages string) error {
	start := time.Now()

	fmt.Println("\n\nRunning Weapon Detection ....")
	err := exec.Command("python3", "./libpy/weapon.py", wepimages).Run()

	fmt.Printf("\nWeapon Detection Time: %v\n", time.Since(start))
	return err
}
