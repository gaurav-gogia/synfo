package lib

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

const (
	linux = "linux"
	mac   = "darwin"
)

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

func PyIdentify(poitest, poitrain, modeltype string) error {
	return exec.Command("python3", "./libpy/face.py", poitest, poitrain, modeltype).Run()
}

func PyDetect(wepimages string) error {
	return exec.Command("python3", "./libpy/weapon.py", wepimages).Run()
}
