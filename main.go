package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"./lib"
)

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
	cli := NewCli()
	if cli.CmdType == EXTCMD {
		in = menu()
	}

	fmt.Println("BufferSize: ", *cli.BufferSize)
	fmt.Println("Imaging ....")
	start := time.Now()
	if err := Run(cli); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("\nImaging Time: %v\n", time.Since(start))

	fmt.Println("\nCalculating Hashes ....")
	integritycheck(*cli.DST)

	handle(getdata(*cli.DST, cli.EviDir, in))
	if cli.CmdType == AUTOCMD {
		fmt.Println("\n\nRunning face recognition ....")
		start = time.Now()
		if err := pyIdentify(*cli.PoI, "./evidence/images/"); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("\nPoI Identification Time: %v\n", time.Since(start))
	}

	fmt.Println("Done!")
}

func getdata(dst, copydst string, in int) error {
	mntloc, copysrc, err := attach(dst)
	if err != nil {
		return err
	}
	lib.Extract(copysrc, copydst, in)
	return detach(mntloc)
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
