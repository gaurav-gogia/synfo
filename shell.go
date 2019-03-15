package main

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func attach(src string) ([]byte, error) {
	return exec.Command("hdiutil", "attach", src).Output()
}

func detach(name string) ([]byte, error) {
	return exec.Command("hdiutil", "detach", name).Output()
}

func getsize(path string) int64 {
	data, _ := exec.Command("diskutil", "info", path).Output()
	info := strings.Split(string(data), "\n")
	var size int64

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

func getnum(data string) int64 {
	const exp = "[^0-9]+"
	num := regexp.MustCompile(exp)
	txt := num.ReplaceAllString(data, "-")
	size, _ := strconv.ParseInt(strings.Split(txt, "-")[3], 10, 64)
	return size
}
