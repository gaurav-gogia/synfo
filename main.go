package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"./lib"
)

type opts struct {
	src        *string
	dst        *string
	poi        *string
	buffersize *int64
	cmdType    string
	evidir     string
}

const (
	defaultBuffer = 10 * 1024 * 1024
	mountinfoPath = "/proc/self/mountinfo"
	partfile      = ".part"
)

// Global Constants
const (
	AUTOCMD = "AUTO"
	EXTCMD  = "EXTRACT"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	in := lib.IMAGE
	dd := cui()
	if dd.cmdType == EXTCMD {
		in = menu()
	}

	/*
		fmt.Println("Imaging ....")
		start := time.Now()
		handle(dd.run())
		fmt.Printf("\nImaging Time: %v\n", time.Since(start))

		fmt.Println("\nCalculating Hashes ....")
		integritycheck(*dd.dst)
	*/

	count, err := getdata(*dd.dst, dd.evidir, in)
	handle(err)

	if dd.cmdType == AUTOCMD && count >= 0 {
		fmt.Println("\nRunning Face Verification ....")
		handle(lib.Verify(*dd.poi, dd.evidir))
	}

	fmt.Println("\n\nDone!")
}

func getdata(dst, copydst string, in int) (int64, error) {
	out, err := attach(dst)

	mntloc := strings.Fields(string(out))[0]
	copysrc := strings.Fields(string(out))[1]

	count := lib.Extract(copysrc, copydst, in)

	_, err = detach(mntloc)

	return count, err
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

func menu() int {
	var in int
	fmt.Println("What do you wish to extract?")
	fmt.Println("1. Picture Files")
	fmt.Println("2. Video Files")
	fmt.Println("3. Audio Files")
	fmt.Println("4. Archive Files")
	fmt.Print("Make your choice: ")
	fmt.Scanln(&in)
	return in
}
