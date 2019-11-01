package lib

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	mntpoint := filepath.Dir(src) + "/" + genname(6)
	if runtime.GOOS == mac {
		out, err := exec.Command("hdiutil", "attach", "-mountpoint", mntpoint, src).Output()
		if err != nil {
			return "", "", err
		}
		mntloc := strings.Fields(string(out))[0]
		copysrc := strings.Fields(string(out))[1]
		return mntloc, copysrc, nil
	} else if runtime.GOOS == linux {
		os.Mkdir(mntpoint, os.ModePerm)
		err := exec.Command("mount", "-o", "loop", src, mntpoint).Run()
		return mntpoint, mntpoint, err
	}

	return "", "", errors.New("unknown runtime")
}

// Detach function unmounts the attached iso file
func Detach(name string) error {
	if runtime.GOOS == mac {
		return exec.Command("hdiutil", "detach", name).Run()
	} else if runtime.GOOS == linux {
		if err := exec.Command("umount", name).Run(); err != nil {
			return err
		}
		return os.Remove(name)
	}
	return errors.New("unknown runtime")
}

func getsize(path string) (int64, error) {
	if runtime.GOOS == mac {
		data, err := exec.Command("diskutil", "info", path).Output()
		if err != nil {
			return 0, err
		}

		info := strings.Split(string(data), "\n")

		for _, str := range info {
			text := fixspace(str)
			if strings.HasPrefix(text, "Disk Size:") {
				return getnum(text)
			}
		}
	} else if runtime.GOOS == linux {
		data, _ := exec.Command("lsblk", "--bytes", path).Output()
		info := strings.Split(string(data), "\n")
		words := strings.Split(string(fixspace(info[1])), " ")
		return strconv.ParseInt(words[3], 10, 64)
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

func getnum(data string) (int64, error) {
	const exp = "[^0-9]+"
	num := regexp.MustCompile(exp)
	txt := num.ReplaceAllString(data, "-")
	return strconv.ParseInt(strings.Split(txt, "-")[3], 10, 64)
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
