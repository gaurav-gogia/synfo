package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CommandLine is the data structure for CLI
type CommandLine struct {
	SRC        *string
	DST        *string
	PoI        *string
	BufferSize *uint64
	CmdType    string
	EviDir     string
	ModelType  *string
}

func (cli *CommandLine) usage() {
	fmt.Println("Usage: ")
	fmt.Println("  auto -src <src_device_file> -dst <dst_file_name> -poi <poi_image_dir> [-buff <buffer_size>] [-model <hog | cnn>] | performs auto forensic image analysis for face verification")
	fmt.Println("  ext -src <src_device_file> -dst <dst_file_name> [-buff <buffer_size>] | images device & extracts specified type of files")
	fmt.Println("  using a fully qualified dir path is recommended for <src_device_file>")
	fmt.Println("  default buffer size is ", defaultBuffer, " bytes")

	fmt.Printf("\n  ---------- EXAMPLE 1 ----------\n")
	fmt.Printf("  ./synfo auto -src /dev/somefile -dst ./somefolder/evi.iso -poi ./person1/images/\n\n")
	fmt.Println("  ---------- EXAMPLE 2 ----------")
	fmt.Printf("  ./synfo auto -src /dev/somefile -dst ./somefolder/evi.iso -poi ./person1/images/ -buff 50000000 -model cnn\n\n")

	fmt.Println("  ---------- EXAMPLE 3 ----------")
	fmt.Printf("  ./synfo ext -src /dev/somefile -dst ./somefolder/evi.iso\n\n")
	fmt.Println("  ---------- EXAMPLE 4 ----------")
	fmt.Printf("  ./synfo ext -src /dev/somefile -dst ./somefolder/evi.iso -buff 50000000\n\n")
}

func (cli *CommandLine) validate() {
	if len(os.Args) < 2 {
		cli.usage()
		os.Exit(0)
	}
}

// NewCli function creates new instances of CLI
func NewCli() CommandLine {
	var cli CommandLine

	cli.validate()

	autocmd := flag.NewFlagSet("auto", flag.ExitOnError)
	extcmd := flag.NewFlagSet("ext", flag.ExitOnError)

	switch os.Args[1] {
	case "auto":
		cli.SRC = autocmd.String("src", "", "Source root directory from where you wish to start scanning")
		cli.DST = autocmd.String("dst", "", "Destination directory where you wish to save output file")
		cli.PoI = autocmd.String("poi", "", "Image directory of suspects' face")
		cli.BufferSize = autocmd.Uint64("buff", defaultBuffer, "Buffer size that you wish to use")
		cli.ModelType = autocmd.String("model", defaultModel, "ML/DL Model to be used")
		handle(autocmd.Parse(os.Args[2:]))
	case "ext":
		cli.SRC = extcmd.String("src", "", "Source root directory from where you wish to start scanning")
		cli.DST = extcmd.String("dst", "", "Destination directory where you wish to save output file")
		cli.BufferSize = extcmd.Uint64("buff", defaultBuffer, "Buffer size that you wish to use")
		handle(extcmd.Parse(os.Args[2:]))
	default:
		cli.usage()
		os.Exit(0)
	}

	if autocmd.Parsed() {
		if *cli.SRC == "" || *cli.DST == "" || *cli.PoI == "" {
			cli.usage()
			os.Exit(0)
		} else if strings.HasSuffix(*cli.SRC, "/") || !(strings.HasPrefix(*cli.SRC, "/dev/")) {
			cli.usage()
			os.Exit(0)
		} else if strings.HasSuffix(*cli.DST, "/") {
			*cli.DST = *cli.DST + "evi.iso"
		} else if !(strings.HasSuffix(*cli.PoI, "/")) {
			*cli.PoI += "/"
		} else if err := sanityCheck(*cli.DST); err != nil {
			cli.usage()
			os.Exit(0)
		} else if *cli.ModelType != "cnn" || *cli.ModelType != "hog" {
			cli.usage()
			os.Exit(0)
		}
		cli.CmdType = AUTOCMD
	}

	if extcmd.Parsed() {
		if *cli.SRC == "" || *cli.DST == "" {
			cli.usage()
			os.Exit(0)
		} else if strings.HasSuffix(*cli.SRC, "/") || !(strings.HasPrefix(*cli.SRC, "/dev/")) {
			cli.usage()
			os.Exit(0)
		} else if strings.HasSuffix(*cli.DST, "/") {
			*cli.DST = *cli.DST + "evi.iso"
		} else if err := sanityCheck(*cli.DST); err != nil {
			cli.usage()
			os.Exit(0)
		}
		cli.CmdType = EXTCMD
	}

	cli.EviDir, _ = filepath.Split(*cli.DST)

	dir, name := filepath.Split(*cli.DST)
	name = strings.TrimSuffix(name, filepath.Ext(name))
	*cli.DST = dir + name + ".iso"

	*cli.BufferSize = fixbuffsize(*cli.BufferSize)

	return cli
}

func fixbuffsize(buffsize uint64) uint64 {
	if buffsize < 50*1024 {
		return buffsize
	}

	fmt.Println("Requested buffersize was too high. Resetting it to: ", defaultBuffer)

	return defaultBuffer
}
