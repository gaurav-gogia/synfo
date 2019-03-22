package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type commandline struct{}

func (cli *commandline) usage() {
	fmt.Println("Usage: ")
	fmt.Println("  auto -src <src_device_file> -dst <dst_file_name> [-buff <buffer_size>] | tries to match keys of all possible types")
	fmt.Println("  using a fully qualified dir path is recommended for <src_device_file>")
	fmt.Println("  default buffer size is ", defaultBuffer)
}

func (cli *commandline) validate() {
	if len(os.Args) < 2 {
		cli.usage()
		os.Exit(0)
	}
}

func cui() opts {
	var dd opts
	var cli commandline

	cli.validate()

	autocmd := flag.NewFlagSet("auto", flag.ExitOnError)

	dd.src = autocmd.String("src", "", "Source root directory from where you wish to start scanning")
	dd.dst = autocmd.String("dst", "", "Destination directory where you wish to save output file")
	dd.buffersize = autocmd.Int64("buff", defaultBuffer, "Buffer size that you wise to provide")

	switch os.Args[1] {
	case "auto":
		handle(autocmd.Parse(os.Args[2:]))
	default:
		cli.usage()
		os.Exit(0)
	}

	if autocmd.Parsed() {
		if *dd.src == "" || *dd.dst == "" {
			cli.usage()
			os.Exit(0)
		} else if strings.HasSuffix(*dd.src, "/") || strings.HasSuffix(*dd.dst, "/") {
			cli.usage()
			os.Exit(0)
		} else if err := sanityCheck(*dd.dst); err != nil {
			cli.usage()
			os.Exit(0)
		}
	}

	return dd
}
