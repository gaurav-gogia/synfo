package main

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"./lib"
	"github.com/cheggaaa/pb"
	"github.com/ulikunitz/xz"
)

func welcome() uint8 {
	var in uint8

	fmt.Println("What do you wish to extract?")

	fmt.Println("1. Image")
	fmt.Println("2. Video")
	fmt.Println("3. Audio")
	fmt.Println("4. Archive")

	fmt.Printf("\nMake your choice: ")
	fmt.Scanf("%d", &in)

	return in
}

func run(dd opts) error {
	if dd.buffer == 0 {
		dd.buffer = defaultBufferSize
	}
	src, size, err := open(dd)
	if err != nil {
		return err
	}
	defer src.Close()
	if err := sanityCheck(dd.dst); err != nil {
		return err
	}
	dst, err := create(dd)
	if err != nil {
		return err
	}
	defer func() {
		dst.Sync()
		dst.Close()
	}()

	w := NewCustomBuff(dst, dd.buffer)
	bar := pb.New64(size).SetUnits(pb.U_BYTES)

	bar.Start()

	_, err = io.Copy(w, bar.NewProxyReader(src))

	bar.Finish()
	return err
}

func open(dd opts) (r io.ReadCloser, size int64, err error) {
	if dd.src == "-" {
		return stdin, 0, nil
	}

	if strings.HasPrefix(dd.src, "http://") || strings.HasPrefix(dd.src, "https://") {
		res, err := http.Get(dd.src)
		if err != nil {
			return nil, 0, err
		}
		size = res.ContentLength
		if size < 0 {
			size = 0
		}
		r = res.Body
	} else {
		r, err = os.Open(dd.src)
		if err != nil {
			return nil, 0, err
		}
	}
	comp := dd.compressionType
	if comp == auto {
		comp = guess(dd.src)
	}

	switch comp {
	case none:
		return r, size, nil
	case gunzip:
		gzr, err := gzip.NewReader(r)
		return gzr, size, err
	case bunzip:
		return ioutil.NopCloser(bzip2.NewReader(r)), size, nil
	case xunzip:
		cr, err := xz.NewReader(r)
		return ioutil.NopCloser(cr), size, err
	}

	panic("can't happen")
}

func guess(src string) uint8 {
	switch filepath.Ext(src) {
	case ".gz":
		return gunzip
	case ".bz2":
		return bunzip
	case ".xz":
		return xunzip
	default:
		return none
	}
}

func sanityCheck(dst string) error {
	f, err := os.Open(mountinfoPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()

	resloved, err := filepath.EvalSymlinks(dst)
	if err == nil {
		dst = resloved
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		l := strings.Fields(scanner.Text())
		if len(l) == 0 {
			continue
		}
		mountPoint := l[4]
		mountSrc := l[9]

		resolvedMountSrc, err := filepath.EvalSymlinks(mountSrc)
		if err == nil {
			mountSrc = resolvedMountSrc
		}

		if strings.HasPrefix(mountSrc, dst) {
			return fmt.Errorf("%s is mounted on %s", mountSrc, mountPoint)
		}
	}
	return scanner.Err()
}

func create(dd opts) (*os.File, error) {
	if dd.dst == "-" {
		return stdout, nil
	}
	return os.Create(dd.dst)
}

// CustomBuff is here because
// go does not offer support to customize the buffer size for
// io.Copy directly, so we need to implement a custom type with:
// ReadFrom and Write
type CustomBuff struct {
	w   io.Writer
	buf []byte
}

// NewCustomBuff is a function
func NewCustomBuff(w io.Writer, size int64) *CustomBuff {
	return &CustomBuff{
		w:   w,
		buf: make([]byte, size),
	}
}

// Write is a function
func (f *CustomBuff) Write(data []byte) (int, error) {
	return f.w.Write(data)
}

func parseargs(args []string) (opts, error) {
	var o opts
	if len(args) == 1 {
		devices, err := findNonCdromRemovableDeviceFiles()
		if err != nil {
			return o, err
		}
		fmt.Printf(`
No target selected, detected the following removable device:
  %s
`, strings.Join(devices, "\n  "))
		return o, fmt.Errorf("please select target device")
	}

	if len(args) == 2 &&
		!strings.Contains(args[0], "=") &&
		!strings.Contains(args[1], "=") {
		o.src = args[0]
		o.dst = args[1]
		return o, nil
	}

	opts := opts{
		buffer: defaultBufferSize,
	}
	for _, arg := range args {
		l := strings.SplitN(arg, "=", 2)
		switch l[0] {
		case "if":
			opts.src = l[1]
		case "of":
			opts.dst = l[1]
		case "bs":
			bs, err := ddAtoi(l[1])
			if err != nil {
				return o, err
			}
			opts.buffer = bs
		case "comp":
			comp, err := ddComp(l[1])
			if err != nil {
				return o, err
			}
			opts.compressionType = comp
		default:
			return o, fmt.Errorf("unknown argument %q", arg)
		}
	}

	return opts, nil
}

func findNonCdromRemovableDeviceFiles() (res []string, err error) {
	devices, err := lib.QueryBySubsystem("block")
	if err != nil {
		return nil, err
	}
	for _, d := range devices {
		if d.GetSysfsAttr("removable") == "1" && d.GetProperty("ID_CDROM") != "1" {
			res = append(res, d.GetDeviceFile())
		}
	}

	return res, nil
}

func ddAtoi(s string) (int64, error) {
	if len(s) < 2 {
		return strconv.ParseInt(s, 10, 64)
	}

	// dd supports suffixes via two chars like "kB"
	fac := int64(1)
	switch s[len(s)-2:] {
	case "kB":
		fac = 1000
	case "MB":
		fac = 1000 * 1000
	case "GB":
		fac = 1000 * 1000 * 1000
	case "TB":
		fac = 1000 * 1000 * 1000 * 1000
	}
	// adjust string if its from xB group
	if fac%10 == 0 {
		s = s[:len(s)-2]
	}

	// check for single char digests
	switch s[len(s)-1] {
	case 'b':
		fac = 512
	case 'K':
		fac = 1024
	case 'M':
		fac = 1024 * 1024
	case 'G':
		fac = 1024 * 1024 * 1024
	case 'T':
		fac = 1024 * 1024 * 1024 * 1024
	}
	// ajust string if its from the X group
	if fac%512 == 0 {
		s = s[:len(s)-1]
	}

	n, err := strconv.ParseInt(s, 10, 64)
	n *= fac
	return n, err
}

func ddComp(s string) (uint8, error) {
	switch s {
	case "auto":
		return auto, nil
	case "none":
		return none, nil
	case "gz", "gzip":
		return gunzip, nil
	case "bz2", "bzip2":
		return bunzip, nil
	case "xz":
		return xunzip, nil
	default:
		return auto, fmt.Errorf("unknown compression type %q", s)
	}
}

func gethashes(imgpath string) (string, string, error) {
	file, err := os.Open(imgpath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", "", err
	}
	hashmd5 := hex.EncodeToString(hash.Sum(nil))

	hash = sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", "", err
	}
	hasha256 := hex.EncodeToString(hash.Sum(nil))

	return hashmd5, hasha256, nil
}
