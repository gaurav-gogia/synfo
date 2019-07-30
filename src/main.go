package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"./cli/"
	"./lib"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	in := 1
	dd := cli.NewCli()
	if dd.cmdType == cli.EXTCMD {
		in = menu()
	}

	fmt.Println("Imaging ....")
	start := time.Now()
	lib.Handle(dd.run())
	fmt.Printf("\nImaging Time: %v\n", time.Since(start))

	fmt.Println("\nCalculating Hashes ....")
	integritycheck(*dd.dst)

	count, err := getdata(*dd.dst, dd.evidir, in)
	lib.Handle(err)

	/*
		TODO:
			Make RPC call or IPC call or something to python script to implement face recognition
	*/

	fmt.Println("\n\nDone!")
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
