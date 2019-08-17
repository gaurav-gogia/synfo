package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"./lib"
)

type opts struct {
	src        *string
	dst        *string
	poi        *string
	buffersize *uint64
	cmdType    string
	evidir     string
}

const (
	defaultBuffer = 10 * 1024
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
	in := 1
	dd := cui()
	if dd.cmdType == EXTCMD {
		in = menu()
	}

	fmt.Println("Imaging ....")
	start := time.Now()
	handle(dd.run())
	fmt.Printf("\nImaging Time: %v\n", time.Since(start))

	fmt.Println("\nCalculating Hashes ....")
	integritycheck(*dd.dst)

	_, err := getdata(*dd.dst, dd.evidir, in)
	handle(err)

	if in == 1 {
		fmt.Println("Running face recognition ....")
		pyIdentify(*dd.poi, dd.evidir+"images/")
	}
}

func getdata(dst, copydst string, in int) (int64, error) {
	mntloc, copysrc, err := attach(dst)
	count := lib.Extract(copysrc, copydst, in)
	err = detach(mntloc)
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
