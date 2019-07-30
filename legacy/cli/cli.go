package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"./lib"
	"github.com/shirou/gopsutil/mem"
)

// CommandLine is the exported data structure for CLI options
type CommandLine struct {
	SRC        *string
	DST        *string
	PoI        *string
	BufferSize *uint64
	CmdType    string
	EviDir     string
}

func (cli *CommandLine) usage() {
	fmt.Println("Usage: ")
	fmt.Println("  auto -src <src_device_file> -dst <dst_file_name> -poi <poi_image_dir> [-buff <buffer_size>] | performs auto forensic image analysis for face verification")
	fmt.Println("  extract -src <src_device_file> -dst <dst_file_name> [-buff <buffer_size>] | images device & extracts specified type of files")
	fmt.Println("  using a fully qualified dir path is recommended for <src_device_file>")
	fmt.Println("  default buffer size is ", defaultBuffer, " bytes")

	fmt.Printf("\n  ---------- EXAMPLE 1 ----------\n")
	fmt.Printf("  ./synfo auto --src /dev/somefile -dst ./somefolder/evi.iso -poi ./person1/images/\n\n")
	fmt.Println("  ---------- EXAMPLE 2 ----------")
	fmt.Printf("  ./synfo auto --src /dev/somefile -dst ./somefolder/evi.iso -poi ./person1/images/ -buff 50000000\n\n")

	fmt.Println("  ---------- EXAMPLE 3 ----------")
	fmt.Printf("  ./synfo extract --src /dev/somefile -dst ./somefolder/evi.iso\n\n")
	fmt.Println("  ---------- EXAMPLE 4 ----------")
	fmt.Printf("  ./synfo extract --src /dev/somefile -dst ./somefolder/evi.iso -buff 50000000\n\n")
}

func (cli *CommandLine) validate() {
	if len(os.Args) < 2 {
		cli.usage()
		os.Exit(0)
	}
}

// NewCli is the entry point function for initializing CLI
func NewCli() CommandLine {
	var cli CommandLine

	cli.validate()

	autocmd := flag.NewFlagSet("auto", flag.ExitOnError)
	extcmd := flag.NewFlagSet("extract", flag.ExitOnError)

	switch os.Args[1] {
	case "auto":
		cli.src = autocmd.String("src", "", "Source root directory from where you wish to start scanning")
		cli.dst = autocmd.String("dst", "", "Destination directory where you wish to save output file")
		cli.poi = autocmd.String("poi", "", "Image directory of suspects' face")
		cli.buffersize = autocmd.Uint64("buff", defaultBuffer, "Buffer size that you wish to use")
		lib.Handle(autocmd.Parse(os.Args[2:]))
	case "extract":
		cli.src = extcmd.String("src", "", "Source root directory from where you wish to start scanning")
		cli.dst = extcmd.String("dst", "", "Destination directory where you wish to save output file")
		cli.buffersize = extcmd.Uint64("buff", defaultBuffer, "Buffer size that you wish to use")
		lib.Handle(extcmd.Parse(os.Args[2:]))
	default:
		cli.usage()
		os.Exit(0)
	}

	if autocmd.Parsed() {
		if *cli.src == "" || *cli.dst == "" || *cli.poi == "" {
			cli.usage()
			os.Exit(0)
		} else if strings.HasSuffix(*cli.src, "/") || !(strings.HasPrefix(*cli.src, "/dev/")) {
			cli.usage()
			os.Exit(0)
		} else if strings.HasSuffix(*cli.dst, "/") {
			*cli.dst = "./evidence/evi.iso"
		} else if !(strings.HasSuffix(*cli.poi, "/")) {
			*cli.poi += "/"
		} else if err := sanityCheck(*cli.dst); err != nil {
			cli.usage()
			os.Exit(0)
		}
		cli.cmdType = AUTOCMD
	}

	if extcmd.Parsed() {
		if *cli.src == "" || *cli.dst == "" {
			cli.usage()
			os.Exit(0)
		} else if strings.HasSuffix(*cli.src, "/") || strings.HasSuffix(*cli.dst, "/") {
			cli.usage()
			os.Exit(0)
		} else if err := sanityCheck(*cli.dst); err != nil {
			cli.usage()
			os.Exit(0)
		}
		cli.cmdType = EXTCMD
	}

	cli.evidir, _ = filepath.Split(*cli.dst)

	dir, name := filepath.Split(*cli.dst)
	name = strings.TrimSuffix(name, filepath.Ext(name))
	*cli.dst = dir + name + ".iso"

	*cli.buffersize = fixbuffsize(*cli.buffersize)

	return cli
}

func fixbuffsize(buffsize uint64) uint64 {
	mem, _ := mem.VirtualMemory()
	newbuff := mem.Free / 4

	if buffsize < newbuff {
		return buffsize
	}

	fmt.Println("Requested buffersize was too high. Resetting it to: ", newbuff)

	return newbuff
}
