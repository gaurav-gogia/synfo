package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"synfo/lib"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	in := 1

	cli, err := lib.NewCli()
	handle(err)

	handle(lib.Run(cli))
	if cli.CmdType == lib.EXTCMD {
		in = menu()
	}
	handle(getdata(*cli.DST, cli.EviDir, in))

	switch cli.CmdType {
	case lib.APDCMD:
		handle(lib.PyApd(cli.EviDir+"images/", *cli.PoI, *cli.ModelType))
	case lib.AWDCMD:
		handle(lib.PyAwd(cli.EviDir + "images/"))
	}

	fmt.Println("\nDone!")
}

func getdata(dst, copydst string, in int) error {
	start := time.Now()
	fmt.Println("\nExtracting Data ....")

	mntloc, copysrc, err := lib.Attach(dst)
	if err != nil {
		return err
	}
	lib.Extract(copysrc, copydst, in)
	err = lib.Detach(mntloc)

	fmt.Printf("\nData Extraction Time: %v\n", time.Since(start))
	return err
}

func handle(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
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
