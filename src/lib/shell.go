package lib

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func attach(src string) (string, string, error) {
	out, err := exec.Command("hdiutil", "attach", src).Output()
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
		if strings.HasPrefix(text, "Disk") {
			size = getnum(text)
			break
		}
	}

	return size
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
