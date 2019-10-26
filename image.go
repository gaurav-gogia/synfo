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

// Run function is the entrypoint for disk imaging, it runs disk imaging
func Run(cli CommandLine) error {
	size, err := getsize(*cli.SRC)
	if err != nil {
		return err
	}

	destination, err := create(*cli.DST)
	if err != nil {
		return err
	}
	destination.Close()

	read, err := unix.Open(*cli.SRC, unix.O_RDONLY, 0777)
	defer unix.Close(read)
	if err != nil {
		return err
	}
	write, err := unix.Open(*cli.DST, unix.O_WRONLY, 0777)
	defer unix.Close(write)
	if err != nil {
		return err
	}

	for i := uint64(0); i < size; i += *cli.BufferSize {
		if size-i <= *cli.BufferSize {
			if err := clone(size-i, read, write); err != nil {
				return err
			}
		} else {
			if err := clone(*cli.BufferSize, read, write); err != nil {
				return err
			}
		}
		fmt.Printf("\rProgress .... %d of %d done", i, size)
	}
	fmt.Println()

	return nil
}

func clone(buffersize uint64, read, write int) error {
	buff := make([]byte, buffersize)

	if _, err := unix.Read(read, buff); err != nil {
		return err
	}

	if _, err := unix.Write(write, buff); err != nil {
		return err
	}

	return nil
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
