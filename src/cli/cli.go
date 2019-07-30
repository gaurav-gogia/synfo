package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shirou/gopsutil/mem"
)

type commandline struct{}

type opts struct {
	src        *string
	dst        *string
	poi        *string
	buffersize *uint64
	cmdType    string
	evidir     string
}

func (cli *commandline) usage() {
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
	extcmd := flag.NewFlagSet("extract", flag.ExitOnError)

	switch os.Args[1] {
	case "auto":
		dd.src = autocmd.String("src", "", "Source root directory from where you wish to start scanning")
		dd.dst = autocmd.String("dst", "", "Destination directory where you wish to save output file")
		dd.poi = autocmd.String("poi", "", "Image directory of suspects' face")
		dd.buffersize = autocmd.Uint64("buff", defaultBuffer, "Buffer size that you wish to use")
		handle(autocmd.Parse(os.Args[2:]))
	case "extract":
		dd.src = extcmd.String("src", "", "Source root directory from where you wish to start scanning")
		dd.dst = extcmd.String("dst", "", "Destination directory where you wish to save output file")
		dd.buffersize = extcmd.Uint64("buff", defaultBuffer, "Buffer size that you wish to use")
		handle(extcmd.Parse(os.Args[2:]))
	default:
		cli.usage()
		os.Exit(0)
	}

	if autocmd.Parsed() {
		if *dd.src == "" || *dd.dst == "" || *dd.poi == "" {
			cli.usage()
			os.Exit(0)
		} else if strings.HasSuffix(*dd.src, "/") || !(strings.HasPrefix(*dd.src, "/dev/")) {
			cli.usage()
			os.Exit(0)
		} else if strings.HasSuffix(*dd.dst, "/") {
			*dd.dst = "./evidence/evi.iso"
		} else if !(strings.HasSuffix(*dd.poi, "/")) {
			*dd.poi += "/"
		} else if err := sanityCheck(*dd.dst); err != nil {
			cli.usage()
			os.Exit(0)
		}
		dd.cmdType = AUTOCMD
	}

	if extcmd.Parsed() {
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
		dd.cmdType = EXTCMD
	}

	dd.evidir, _ = filepath.Split(*dd.dst)

	dir, name := filepath.Split(*dd.dst)
	name = strings.TrimSuffix(name, filepath.Ext(name))
	*dd.dst = dir + name + ".iso"

	*dd.buffersize = fixbuffsize(*dd.buffersize)

	return dd
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
