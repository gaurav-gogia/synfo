package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"./lib"
)

type opts struct {
	src           *string
	dst           *string
	buffersize    *int64
	maxthreads    *int
	multithreaded *bool
}

const (
	defaultBuffer = 10 * 1024 * 1024
	mountinfoPath = "/proc/self/mountinfo"
	partfile      = ".part"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	in := menu()
	dd := cui()

	fmt.Println("Processing....")
	start := time.Now()
	handle(dd.run())
	fmt.Printf("\nImaging Time: %v\n", time.Since(start))

	fmt.Println("\nCalculating Hashes ....")
	integritycheck(*dd.dst)

	handle(getdata(*dd.dst, in))

	fmt.Println("Done!")
}

func getdata(dst string, in int8) error {
	out, err := attach(dst)

	mntloc := strings.Fields(string(out))[0]
	copysrc := strings.Fields(string(out))[1]

	lib.Extract(copysrc, dst, in)

	_, err = detach(mntloc)

	return err
}

func integritycheck(dst string) {
	md, sha, err := gethashes(dst)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to gain hashes: %v", err))
		os.Exit(1)
	}
	fmt.Println("MD5: ", md)
	fmt.Println("SHA256: ", sha)
}

func handle(err error) {
	if err != nil {
		fmt.Println(err)
	}
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
