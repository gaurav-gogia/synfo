package main

import (
	"crypto/rand"
	"encoding/base32"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func attach(src string) (string, string, error) {
	mntpoint := genname(6)
	out, err := exec.Command("hdiutil", "attach", "-mountpoint", mntpoint, src).Output()
	if err != nil {
		return "", "", err
	}
	mntloc := strings.Fields(string(out))[0]
	copysrc := strings.Fields(string(out))[1]
	return mntloc, copysrc, nil
}

func detach(name string) error {
	_, err := exec.Command("hdiutil", "detach", name).Output()
	return err
}

func getsize(path string) uint64 {
	data, _ := exec.Command("diskutil", "info", path).Output()
	info := strings.Split(string(data), "\n")
	var size uint64

	for _, str := range info {
		text := fixspace(str)
		if strings.HasPrefix(text, "Disk Size:") {
			size = getnum(text)
			break
		}
	}

	return size
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

func pyIdentify(poitest, poitrain, modeltype string) error {
	return exec.Command("python3", "./libpy/face.py", poitest, poitrain, modeltype).Run()
}
