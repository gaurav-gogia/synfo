package lib

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
	"time"
)

// Clone function is the entrypoint for disk imaging, it runs disk imaging
func Clone(cli CommandLine) error {
	start := time.Now()

	if !confirm(cli.DST) {
		return nil
	}

	fmt.Println("\nImaging ....")
	fmt.Println("Using BufferSize: ", cli.BufferSize)

	read, write, size, err := setup(cli.SRC, cli.DST)
	if err != nil {
		return err
	}
	defer read.Close()
	defer write.Close()

	for i := int64(0); i < size; i += cli.BufferSize {
		percent := (float64(i) / float64(size)) * 100
		fmt.Printf("\rProgress .... %f%%", percent)

		if i == size {
			fmt.Printf("\rProgress .... %f%%", percent)
			break
		}

		if size-i <= cli.BufferSize {
			cli.BufferSize = size - i
		}

		if err := copyData(cli.BufferSize, read, write); err != nil {
			return err
		}
	}
	fmt.Println()

	fmt.Printf("\nImaging Time: %v\n", time.Since(start))
	integritycheck(cli.DST)

	return nil
}

func setup(src, dst string) (*os.File, *os.File, int64, error) {
	size, err := getsize(src)
	if err != nil {
		return nil, nil, 0, err
	}

	var read, write *os.File
	if read, err = os.Open(src); err != nil {
		return nil, nil, 0, err
	}
	if write, err = os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, 0600); err != nil {
		return nil, nil, 0, err
	}

	return read, write, size, nil
}

func copyData(buffersize int64, read, write *os.File) error {
	buff := make([]byte, buffersize)

	if _, err := read.Read(buff); err != nil {
		return err
	}

	if _, err := write.Write(buff); err != nil {
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

func integritycheck(dst string) {
	start := time.Now()
	fmt.Println("\nCalculating Hashes ....")

	md, sha, err := gethashes(dst)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to gain hashes: %v", err))
		os.Exit(1)
	}
	fmt.Println("MD5: ", md)
	fmt.Println("SHA256: ", sha)

	fmt.Printf("Hash Calculation Time: %v\n", time.Since(start))
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

func confirm(dst string) bool {
	var ans string
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		return true
	}

	fmt.Printf("Disk Image already exists. Continue [y/n]? -> ")
	fmt.Scanln(&ans)

	if ans == "y" || ans == "Y" {
		return true
	}

	fmt.Println("User skipped disk imaging....")
	return false
}
