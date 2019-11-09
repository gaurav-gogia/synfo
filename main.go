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
	cli, err := lib.NewCli()
	handle(err)

	handle(lib.Run(cli))
	handle(getdata(cli.DST, cli.EviDir, cli.FileType))

	switch cli.CmdType {
	case lib.APDCMD:
		handle(lib.PyApd(cli.EviDir+"images/", cli.PoI, cli.ModelType))
	case lib.AWDCMD:
		handle(lib.PyAwd(cli.EviDir + "images/"))
	}

	fmt.Println("\nDone!")
}

func getdata(dst, copydst string, ft string) error {
	start := time.Now()
	fmt.Println("\nExtracting Data ....")

	mntloc, copysrc, err := lib.Attach(dst)
	if err != nil {
		return err
	}
	lib.Extract(copysrc, copydst, ft)
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
