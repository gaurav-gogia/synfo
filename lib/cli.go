package lib

import (
	"errors"
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
	BufferSize *int64
	CmdType    string
	EviDir     string
	ModelType  *string
	FileType   *string
	Help       *bool
	Examples   *bool
}

func (cli *CommandLine) usage() {
	fmt.Println("synfo - the sharpest sentry, in your arsenal.")
	fmt.Println("synfo is an automated digital forensic investigation framework")

	fmt.Printf("\n\nUSAGE: synfo COMMAND [FLAGS...]")

	fmt.Printf("\n\nFLAGS:")
	fmt.Printf("\n -h, --help")
	fmt.Printf("\n\t%s", helpusageflag)
	fmt.Printf("\n -e, --examples")
	fmt.Printf("\n\t%s", exampleusageflag)

	fmt.Printf("\n\nCommands:")
	fmt.Printf("\n  %s  ->  %s", EXTCMD, extcmduse)
	fmt.Printf("\n  %s  ->  %s", APDCMD, apdcmduse)
	fmt.Printf("\n  %s  ->  %s", AWDCMD, awdcmduse)

	os.Exit(0)
}

func (cli *CommandLine) examples() {
	extexamples()
	apdexamples()
	awdexamples()
	os.Exit(0)
}

func extexamples() {
	fmt.Printf("\n\nExamples: ext\n\n")
	fmt.Println("  ---------- EXAMPLE 1 ----------")
	fmt.Printf("  ./synfo ext -src /dev/somefile -dst ./somefolder/evi.iso\n\n")
	fmt.Println("  ---------- EXAMPLE 2 ----------")
	fmt.Printf("  ./synfo ext -src /dev/somefile -dst ./somefolder/evi.iso -bs 50000000\n\n")
	fmt.Println("  ---------- EXAMPLE 3 ----------")
	fmt.Printf("  ./synfo ext -src /dev/somefile -dst ./somefolder/evi.iso -ft audio\n\n")
}
func apdexamples() {
	fmt.Printf("\n\nExamples: apd\n\n")
	fmt.Printf("\n  ---------- EXAMPLE 1 ----------\n")
	fmt.Printf("  ./synfo apd -src /dev/somefile -dst ./somefolder/evi.iso -poi ./person1/images/\n\n")
	fmt.Println("  ---------- EXAMPLE 2 ----------")
	fmt.Printf("  ./synfo apd -src /dev/somefile -dst ./somefolder/evi.iso -poi ./person1/images/ -bs 50000000 -model cnn\n\n")
}
func awdexamples() {
	fmt.Printf("\n\nExamples: awd\n\n")
	fmt.Println("  ---------- EXAMPLE 6 ----------")
	fmt.Printf("  ./synfo awd -src /dev/somefile -dst ./somefolder/evi.iso\n\n")
	fmt.Println("  ---------- EXAMPLE 7 ----------")
	fmt.Printf("  ./synfo awd -src /dev/somefile -dst ./somefolder/evi.iso -bs 50000000\n\n")
}

func (cli *CommandLine) extusage() {
	basicusage(EXTCMD, extcmduse)

	fmt.Printf("\n\n -ft [default: %s]", defaultFt)
	fmt.Printf("\n\t%s", ftflaghelp)

	extexamples()

	os.Exit(0)
}
func (cli *CommandLine) apdusage() {
	basicusage(APDCMD, apdcmduse)

	fmt.Printf("\n\n -model [default: %s]", defaultModel)
	fmt.Printf("\n\t%s", modelflaghelp)

	apdexamples()

	os.Exit(0)
}
func (cli *CommandLine) awdusage() {
	basicusage(AWDCMD, awdcmduse)
	awdexamples()
	os.Exit(0)
}

func basicusage(cmdname, cmduse string) {
	fmt.Printf("synfo %s -> %s", cmdname, cmduse)
	fmt.Printf("\n\nUSAGE: synfo %s [FLAGS...]", cmdname)

	fmt.Printf("\n\nFLAGS:")
	fmt.Printf("\n -src")
	fmt.Printf("\n\t%s", srcflaghelp)
	fmt.Printf("\n -dst")
	fmt.Printf("\n\t%s", dstflaghelp)

	fmt.Printf("\n\nOPTIONAL FLAGS:")
	fmt.Printf("\n -bs [default: %d]", defaultBuffer)
	fmt.Printf("\n\t%s", bsflaghelp)
}

func (cli *CommandLine) validate() error {
	switch len(os.Args) {
	case 1:
		cli.usage()
	case 2:
		if (os.Args[1] == "-h") || (os.Args[1] == "--help") {
			cli.usage()
		} else if (os.Args[1] == "-e") || (os.Args[1] == "--examples") {
			cli.examples()
		} else if (os.Args[1] == "ext") || (os.Args[1] == "apd") || (os.Args[1] == "awd") {
			return nil
		} else {
			return errors.New("Unknown argument: " + os.Args[1])
		}
	}
	return nil
}

// NewCli function creates new instances of CLI
func NewCli() (CommandLine, error) {
	var cli CommandLine
	var err error

	if err := cli.validate(); err != nil {
		return cli, err
	}

	apdcmd := flag.NewFlagSet("apd", flag.ExitOnError)
	extcmd := flag.NewFlagSet("ext", flag.ExitOnError)
	awdcmd := flag.NewFlagSet("awd", flag.ExitOnError)

	switch os.Args[1] {
	case EXTCMD:
		cli.SRC = extcmd.String("src", "", srcflaghelp)
		cli.DST = extcmd.String("dst", "", dstflaghelp)
		cli.FileType = extcmd.String("ft", defaultFt, ftflaghelp)
		cli.BufferSize = extcmd.Int64("bs", defaultBuffer, bsflaghelp)
		if err := extcmd.Parse(os.Args[2:]); err != nil {
			return cli, err
		}
	case APDCMD:
		cli.SRC = apdcmd.String("src", "", srcflaghelp)
		cli.DST = apdcmd.String("dst", "", dstflaghelp)
		cli.PoI = apdcmd.String("poi", "", poiflaghelp)
		cli.BufferSize = apdcmd.Int64("bs", defaultBuffer, bsflaghelp)
		cli.ModelType = apdcmd.String("model", defaultModel, modelflaghelp)
		if err := apdcmd.Parse(os.Args[2:]); err != nil {
			return cli, err
		}
	case AWDCMD:
		cli.SRC = awdcmd.String("src", "", srcflaghelp)
		cli.DST = awdcmd.String("dst", "", dstflaghelp)
		cli.BufferSize = extcmd.Int64("bs", defaultBuffer, bsflaghelp)
		if err := awdcmd.Parse(os.Args[2:]); err != nil {
			return cli, err
		}
	default:
		cli.usage()
	}

	if extcmd.Parsed() {
		cli.parseExt()
	}
	if apdcmd.Parsed() {
		cli.parseApd()
	}
	if awdcmd.Parsed() {
		cli.parseAwd()
	}

	err = cli.finetuning()

	return cli, err
}

func fixbuffsize(buffsize int64) int64 {
	if buffsize < 50*1024 {
		return buffsize
	}

	fmt.Println("Requested buffersize was too high. Resetting it to: ", defaultBuffer)

	return defaultBuffer
}

func (cli *CommandLine) parseExt() {
	if len(os.Args) < 4 {
		cli.extusage()
	}

	if *cli.SRC == "" || *cli.DST == "" {
		cli.extusage()
	} else if strings.HasSuffix(*cli.SRC, "/") || !(strings.HasPrefix(*cli.SRC, "/dev/")) {
		cli.extusage()
	} else if strings.HasSuffix(*cli.DST, "/") {
		*cli.DST = *cli.DST + defaultDiskImage
	} else if err := sanityCheck(*cli.DST); err != nil {
		cli.extusage()
	}

	switch *cli.FileType {
	case IMAGE:
		fallthrough
	case AUDIO:
		fallthrough
	case VIDEO:
		fallthrough
	case ARCHIVE:
	default:
		cli.extusage()
	}

	cli.CmdType = EXTCMD
}

func (cli *CommandLine) parseApd() {
	if len(os.Args) < 5 {
		cli.apdusage()
	}

	if *cli.SRC == "" || *cli.DST == "" || *cli.PoI == "" {
		cli.apdusage()
	} else if strings.HasSuffix(*cli.SRC, "/") || !(strings.HasPrefix(*cli.SRC, "/dev/")) {
		cli.apdusage()
	} else if strings.HasSuffix(*cli.DST, "/") {
		*cli.DST = *cli.DST + defaultDiskImage
	} else if !(strings.HasSuffix(*cli.PoI, "/")) {
		*cli.PoI += "/"
	} else if err := sanityCheck(*cli.DST); err != nil {
		cli.apdusage()
	}

	switch *cli.ModelType {
	case "cnn":
		fallthrough
	case "hog":
	default:
		cli.apdusage()
	}

	*cli.FileType = defaultFt
	cli.CmdType = APDCMD
}

func (cli *CommandLine) parseAwd() {
	if len(os.Args) < 4 {
		cli.awdusage()
	}

	if *cli.SRC == "" || *cli.DST == "" {
		cli.awdusage()
	} else if strings.HasSuffix(*cli.SRC, "/") || !(strings.HasPrefix(*cli.SRC, "/dev/")) {
		cli.awdusage()
	} else if strings.HasSuffix(*cli.DST, "/") {
		*cli.DST = *cli.DST + defaultDiskImage
	} else if err := sanityCheck(*cli.DST); err != nil {
		cli.awdusage()
	}

	*cli.FileType = defaultFt
	cli.CmdType = AWDCMD
}

func (cli *CommandLine) finetuning() error {
	var err error

	*cli.SRC, err = filepath.Abs(*cli.SRC)
	if err != nil {
		return errors.New("Could NOT convert" + *cli.SRC + "into Absolute path")
	}
	*cli.DST, err = filepath.Abs(*cli.DST)
	if err != nil {
		return errors.New("Could NOT convert" + *cli.DST + "into Absolute path")
	}

	cli.EviDir, _ = filepath.Split(*cli.DST)

	dir, name := filepath.Split(*cli.DST)
	name = strings.TrimSuffix(name, filepath.Ext(name))
	*cli.DST = dir + name + ".iso"

	*cli.BufferSize = fixbuffsize(*cli.BufferSize)

	return err
}
