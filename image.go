package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix"
)

func (dd *opts) run() error {
	size := getsize(*dd.src)

	for i := int64(0); i < size; i += *dd.buffersize {
		readwrite(*dd.buffersize, *dd.src, *dd.dst)
		fmt.Printf("\rImaging .... %d of %d done", i, size)
	}
	fmt.Println()

	return nil
}

func readwrite(buffersize int64, src, dst string) {
	destination, err := create(dst)
	handle(err)
	destination.Close()

	buff := make([]byte, buffersize)

	read, err := unix.Open(src, unix.O_RDONLY, 0777)
	defer unix.Close(read)
	handle(err)
	write, err := unix.Open(dst, unix.O_WRONLY, 0777)
	defer unix.Close(write)
	handle(err)

	_, err = unix.Read(read, buff)
	handle(err)
	_, err = unix.Write(write, buff)
	handle(err)
}

func sanityCheck(dst string) error {
	f, err := os.Open(mountinfoPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()

	resloved, err := filepath.EvalSymlinks(dst)
	if err == nil {
		dst = resloved
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		l := strings.Fields(scanner.Text())
		if len(l) == 0 {
			continue
		}
		mountPoint := l[4]
		mountSrc := l[9]

		resolvedMountSrc, err := filepath.EvalSymlinks(mountSrc)
		if err == nil {
			mountSrc = resolvedMountSrc
		}

		if strings.HasPrefix(mountSrc, dst) {
			return fmt.Errorf("%s is mounted on %s", mountSrc, mountPoint)
		}
	}
	return scanner.Err()
}

func gethashes(imgpath string) (string, string, error) {
	file, err := os.Open(imgpath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", "", err
	}
	hashmd5 := hex.EncodeToString(hash.Sum(nil))

	hash = sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", "", err
	}
	hasha256 := hex.EncodeToString(hash.Sum(nil))

	return hashmd5, hasha256, nil
}

func create(dst string) (*os.File, error) {
	if dst == "-" {
		return os.Stdout, nil
	}
	return os.Create(dst)
}
