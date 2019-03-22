package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

func (dd *opts) run() error {
	var thread int64
	var names []string
	done := make(chan bool)
	size := getsize(*dd.src)

	for i := int64(0); i < size; i += *dd.buffersize {
		dstname := partfile + strconv.FormatInt(i, 10)
		names = append(names, dstname)

		go parts(*dd.buffersize, i, *dd.src, dstname, done)
		thread++

		if thread > 16 {
			<-done
			thread--
		}

		fmt.Printf("\rSpawning %d threads", thread)
	}
	fmt.Println()

	for thread > 0 {
		<-done
		thread--
	}

	file, err := os.OpenFile(*dd.dst, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	defer file.Close()
	handle(err)

	for _, name := range names {
		data, err := ioutil.ReadFile(name)
		handle(err)

		_, err = file.Write(data)
		handle(err)

		os.Remove(name)
		fmt.Printf("\rCompiling .... %s", name)
	}
	fmt.Println()

	return nil
}

func parts(buffersize, offset int64, src, dst string, done chan bool) {
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

	unix.Seek(read, offset, 0)
	_, err = unix.Read(read, buff)
	handle(err)

	_, err = unix.Write(write, buff)
	handle(err)

	done <- true
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
