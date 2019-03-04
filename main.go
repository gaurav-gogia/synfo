package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"./lib"
)

type opts struct {
	src             string
	dst             string
	buffer          int64
	compressionType uint8
}

var (
	stdin         = os.Stdin
	stdout        = os.Stdout
	mountinfoPath = "/proc/self/mountinfo"
	extOut        = "/Volumes/store/evidence/data/"
)

const (
	none uint8 = 1 << iota
	gunzip
	bunzip
	xunzip
	auto = 0

	defaultBufferSize = int64(10 * 1024 * 1024)
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	in := menu()

	dd, err := parseargs(os.Args[1:])
	if err != nil {
		fmt.Println(fmt.Errorf("failed to parse args: %v", err))
		os.Exit(1)
	}

	if err := run(dd); err != nil {
		fmt.Println(fmt.Errorf("failed to dd: %v", err))
		os.Exit(1)
	}

	md, sha, err := gethashes(dd.dst)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to gain hashes: %v", err))
		os.Exit(1)
	}
	fmt.Println("MD5: ", md)
	fmt.Println("SHA256: ", sha)

	out, err := attach(dd.dst)

	mntloc := strings.Fields(string(out))[0]
	copysrc := strings.Fields(string(out))[1]

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(copysrc)

	lib.Extract(copysrc, dd.dst, in)

	_, err = detach(mntloc)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Done!")
}

func attach(src string) ([]byte, error) {
	return exec.Command("hdiutil", "attach", src).Output()
}

func detach(name string) ([]byte, error) {
	return exec.Command("hdiutil", "detach", name).Output()
}

func menu() int8 {
	var in int8
	fmt.Println("What do you wish to extract?")
	fmt.Println("1. Picture Files")
	fmt.Println("2. Video Files")
	fmt.Println("3. Audio Files")
	fmt.Println("4. Archive Files")
	fmt.Print("Make your choice: ")
	fmt.Scanln(&in)
	return in
}
