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
	SRC        string
	DST        string
	PoI        string
	BufferSize int64
	CmdType    string
	EviDir     string
	ModelType  string
	FileType   string
	Help       bool
	Examples   bool
	Flash      bool
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
	cli.extexamples()
	cli.apdexamples()
	cli.awdexamples()
	os.Exit(0)
}

func (cli *CommandLine) extexamples() {
	fmt.Printf("\n\nExamples: ext\n\n")
	fmt.Println("  ---------- EXAMPLE 1 ----------")
	fmt.Printf("  ./synfo ext -src /dev/somefile -dst ./somefolder/evi.iso\n\n")
	fmt.Println("  ---------- EXAMPLE 2 ----------")
	fmt.Printf("  ./synfo ext -src /dev/somefile -dst ./somefolder/evi.iso -bs 50000000\n\n")
	fmt.Println("  ---------- EXAMPLE 3 ----------")
	fmt.Printf("  ./synfo ext -src /dev/somefile -dst ./somefolder/evi.iso -ft audio\n\n")
}
func (cli *CommandLine) apdexamples() {
	fmt.Printf("\n\nExamples: apd\n\n")
	fmt.Printf("\n  ---------- EXAMPLE 1 ----------\n")
	fmt.Printf("  ./synfo apd -src /dev/somefile -dst ./somefolder/evi.iso -poi ./person1/images/\n\n")
	fmt.Println("  ---------- EXAMPLE 2 ----------")
	fmt.Printf("  ./synfo apd -src /dev/somefile -dst ./somefolder/evi.iso -poi ./person1/images/ -bs 50000000 -model cnn\n\n")
}
func (cli *CommandLine) awdexamples() {
	fmt.Printf("\n\nExamples: awd\n\n")
	fmt.Println("  ---------- EXAMPLE 1 ----------")
	fmt.Printf("  ./synfo awd -src /dev/somefile -dst ./somefolder/evi.iso\n\n")
	fmt.Println("  ---------- EXAMPLE 2 ----------")
	fmt.Printf("  ./synfo awd -src /dev/somefile -dst ./somefolder/evi.iso -bs 50000000\n\n")
}

func (cli *CommandLine) exthelp() {
	basichelp(EXTCMD, extcmduse)

	fmt.Printf("\n\n -ft [default: %s]", defaultFt)
	fmt.Printf("\n\t%s", ftflaghelp)

	os.Exit(0)
}
func (cli *CommandLine) apdhelp() {
	basichelp(APDCMD, apdcmduse)

	fmt.Printf("\n\n -model [default: %s]", defaultModel)
	fmt.Printf("\n\t%s", modelflaghelp)

	os.Exit(0)
}
func (cli *CommandLine) awdhelp() {
	basichelp(AWDCMD, awdcmduse)
	os.Exit(0)
}

func basichelp(cmdname, cmduse string) {
	fmt.Printf("synfo %s -> %s", cmdname, cmduse)
	fmt.Printf("\n\nUSAGE: synfo %s [FLAGS...]", cmdname)

	fmt.Printf("\n\nFLAGS:")
	fmt.Printf("\n -h, --help")
	fmt.Printf("\n\t%s", helpusageflag)
	fmt.Printf("\n -e, --examples")
	fmt.Printf("\n\t%s", exampleusageflag)

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
		extcmd.StringVar(&cli.SRC, "src", "", srcflaghelp)
		extcmd.StringVar(&cli.DST, "dst", "", dstflaghelp)
		extcmd.StringVar(&cli.FileType, "ft", defaultFt, ftflaghelp)
		extcmd.Int64Var(&cli.BufferSize, "bs", defaultBuffer, bsflaghelp)

		extcmd.BoolVar(&cli.Help, "h", false, helpusageflag)
		extcmd.BoolVar(&cli.Help, "help", false, helpusageflag)
		extcmd.BoolVar(&cli.Examples, "e", false, exampleusageflag)
		extcmd.BoolVar(&cli.Examples, "examples", false, exampleusageflag)

		extcmd.BoolVar(&cli.Flash, "flash", false, flashusageflag)
		if err := extcmd.Parse(os.Args[2:]); err != nil {
			return cli, err
		}
	case APDCMD:
		apdcmd.StringVar(&cli.SRC, "src", "", srcflaghelp)
		apdcmd.StringVar(&cli.DST, "dst", "", dstflaghelp)
		apdcmd.Int64Var(&cli.BufferSize, "bs", defaultBuffer, bsflaghelp)
		apdcmd.StringVar(&cli.ModelType, "model", defaultModel, modelflaghelp)

		apdcmd.BoolVar(&cli.Help, "h", false, helpusageflag)
		apdcmd.BoolVar(&cli.Help, "help", false, helpusageflag)
		apdcmd.BoolVar(&cli.Examples, "e", false, exampleusageflag)
		apdcmd.BoolVar(&cli.Examples, "examples", false, exampleusageflag)

		apdcmd.BoolVar(&cli.Flash, "flash", false, flashusageflag)
		if err := apdcmd.Parse(os.Args[2:]); err != nil {
			return cli, err
		}
	case AWDCMD:
		awdcmd.StringVar(&cli.SRC, "src", "", srcflaghelp)
		awdcmd.StringVar(&cli.DST, "dst", "", dstflaghelp)
		awdcmd.Int64Var(&cli.BufferSize, "bs", defaultBuffer, bsflaghelp)

		awdcmd.BoolVar(&cli.Help, "h", false, helpusageflag)
		awdcmd.BoolVar(&cli.Help, "help", false, helpusageflag)
		awdcmd.BoolVar(&cli.Examples, "e", false, exampleusageflag)
		awdcmd.BoolVar(&cli.Examples, "examples", false, exampleusageflag)

		awdcmd.BoolVar(&cli.Flash, "flash", false, flashusageflag)
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
	if cli.Help {
		cli.exthelp()
	}
	if cli.Examples {
		cli.extexamples()
		os.Exit(0)
	}
	if len(os.Args) < 4 {
		cli.exthelp()
	}

	if cli.SRC == "" || cli.DST == "" {
		cli.exthelp()
	} else if strings.HasSuffix(cli.SRC, "/") || !(strings.HasPrefix(cli.SRC, "/dev/")) {
		cli.exthelp()
	} else if strings.HasSuffix(cli.DST, "/") {
		cli.DST = cli.DST + defaultDiskImage
	} else if err := sanityCheck(cli.DST); err != nil {
		cli.exthelp()
	}

	switch cli.FileType {
	case IMAGE:
		fallthrough
	case AUDIO:
		fallthrough
	case VIDEO:
		fallthrough
	case ARCHIVE:
	default:
		cli.exthelp()
	}

	cli.CmdType = EXTCMD
}

func (cli *CommandLine) parseApd() {
	if cli.Help {
		cli.apdhelp()
	}
	if cli.Examples {
		cli.apdexamples()
		os.Exit(0)
	}
	if len(os.Args) < 4 {
		cli.apdhelp()
	}

	if cli.SRC == "" || cli.DST == "" || cli.PoI == "" {
		cli.apdhelp()
	} else if strings.HasSuffix(cli.SRC, "/") || !(strings.HasPrefix(cli.SRC, "/dev/")) {
		cli.apdhelp()
	} else if strings.HasSuffix(cli.DST, "/") {
		cli.DST = cli.DST + defaultDiskImage
	} else if !(strings.HasSuffix(cli.PoI, "/")) {
		cli.PoI += "/"
	} else if err := sanityCheck(cli.DST); err != nil {
		cli.apdhelp()
	}

	switch cli.ModelType {
	case "cnn":
		fallthrough
	case "hog":
	default:
		cli.apdhelp()
	}

	cli.FileType = defaultFt
	cli.CmdType = APDCMD
}

func (cli *CommandLine) parseAwd() {
	if cli.Help {
		cli.awdhelp()
	}
	if cli.Examples {
		cli.awdexamples()
		os.Exit(0)
	}
	if len(os.Args) < 4 {
		cli.awdhelp()
	}

	if cli.SRC == "" || cli.DST == "" {
		cli.awdhelp()
	} else if strings.HasSuffix(cli.SRC, "/") || !(strings.HasPrefix(cli.SRC, "/dev/")) {
		cli.awdhelp()
	} else if strings.HasSuffix(cli.DST, "/") {
		cli.DST = cli.DST + defaultDiskImage
	} else if err := sanityCheck(cli.DST); err != nil {
		cli.awdhelp()
	}

	cli.FileType = defaultFt
	cli.CmdType = AWDCMD
}

func (cli *CommandLine) finetuning() error {
	var err error

	cli.SRC, err = filepath.Abs(cli.SRC)
	if err != nil {
		return errors.New("Could NOT convert" + cli.SRC + "into Absolute path")
	}
	cli.DST, err = filepath.Abs(cli.DST)
	if err != nil {
		return errors.New("Could NOT convert" + cli.DST + "into Absolute path")
	}

	cli.EviDir, _ = filepath.Split(cli.DST)

	dir, name := filepath.Split(cli.DST)
	name = strings.TrimSuffix(name, filepath.Ext(name))
	cli.DST = dir + name + ".iso"

	cli.BufferSize = fixbuffsize(cli.BufferSize)

	return err
}
