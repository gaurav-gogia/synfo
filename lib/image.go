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

	"golang.org/x/sys/unix"
)

// Run function is the entrypoint for disk imaging, it runs disk imaging
func Run(cli CommandLine) error {
	start := time.Now()

	if !confirm(*cli.DST) {
		return nil
	}

	fmt.Println("\nImaging ....")
	fmt.Println("Using BufferSize: ", *cli.BufferSize)

	size, err := getsize(*cli.SRC)
	if err != nil {
		return err
	}

	read, write, err := setup(*cli.SRC, *cli.DST)
	if err != nil {
		return err
	}
	defer unix.Close(read)
	defer unix.Close(write)

	for i := int64(0); i <= size; i += *cli.BufferSize {
		percent := (float64(i) / float64(size)) * 100
		fmt.Printf("\rProgress .... %f%%", percent)

		if i == size {
			break
		}

		if size-i <= *cli.BufferSize {
			if err := clone(size-i, read, write); err != nil {
				return err
			}
		} else {
			if err := clone(*cli.BufferSize, read, write); err != nil {
				return err
			}
		}
	}
	fmt.Println()

	fmt.Printf("\nImaging Time: %v\n", time.Since(start))
	integritycheck(*cli.DST)

	return nil
}

func setup(src, dst string) (int, int, error) {
	destination, err := create(dst)
	if err != nil {
		return 0, 0, err
	}
	destination.Close()

	read, err := unix.Open(src, unix.O_RDONLY, 0777)
	if err != nil {
		return 0, 0, err
	}

	write, err := unix.Open(dst, unix.O_WRONLY, 0777)
	if err != nil {
		return 0, 0, err
	}

	return read, write, nil
}

func clone(buffersize int64, read, write int) error {
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

func create(dst string) (*os.File, error) {
	if dst == "-" {
		return os.Stdout, nil
	}
	return os.Create(dst)
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
