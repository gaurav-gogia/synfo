package lib

import (
	"archive/zip"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/h2non/filetype"
	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdf "github.com/unidoc/unidoc/pdf/model"
	"golang.org/x/sys/unix"
)

// Extract function takes in root dir, destination and choice.
// Walks through entire file structure to copy files based on their magic numbers
func Extract(root, dst string, in int) {
	var count int64
	filepath.Walk(root, func(filepath string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Print("\n", err, "\n\n")
			return nil
		}

		if !info.IsDir() {
			buf, err := ioutil.ReadFile(filepath)
			switch in {
			case IMAGE:
				if err := copyimage(&buf, &count, dst+"images/", info.Name()); err != nil {
					fmt.Println(err)
				}
				if err := getimgfrompdf(&buf, &count, dst+"images/", filepath); err != nil {
					fmt.Println(err)
				}
				if err := getimgfrompop(&buf, &count, dst+"images/", filepath); err != nil {
					fmt.Println(err)
				}
				if err := carvefile(&count, dst+"images/", filepath); err != nil {
					fmt.Println(err)
				}
			case VIDEO:
				if err := copyvideo(&buf, &count, dst+"videos/", info.Name()); err != nil {
					fmt.Println(err)
				}
			case AUDIO:
				if err := copyaudio(&buf, &count, dst+"audios/", info.Name()); err != nil {
					fmt.Println(err)
				}
			case ARCHIVE:
				if err := copyarchive(&buf, &count, dst+"archives/", info.Name()); err != nil {
					fmt.Println(err)
				}
			default:
				fmt.Println("Wrong Choice")
				return nil
			}

			return err
		}
		return nil
	})

	fmt.Println("\nTotal files found: ", count)
}

func getimgfrompdf(buf *[]byte, count *int64, dst, infile string) error {
	filetypedir(dst)

	if ispdf(*buf) {
		f, err := os.Open(infile)
		defer f.Close()
		if err != nil {
			return err
		}

		pdfreader, err := pdf.NewPdfReader(f)
		if err != nil {
			return err
		}

		isenc, err := pdfreader.IsEncrypted()
		if err != nil {
			return err
		}

		if isenc {
			return errors.New("File is encrypted")
		}

		numPages, err := pdfreader.GetNumPages()
		if err != nil {
			return err
		}

		for i := 0; i < numPages; i++ {
			page, err := pdfreader.GetPage(i + 1)
			if err != nil {
				return err
			}

			if err := extractImagesOnPage(page, dst, infile, count); err != nil {
				return err
			}
		}
	}

	return nil
}

func extractImagesOnPage(page *pdf.PdfPage, dst, infile string, count *int64) error {
	contents, err := page.GetAllContentStreams()
	if err != nil {
		return err
	}

	return extImgsInContStream(contents, page.Resources, dst, infile, count)
}

func extImgsInContStream(contents string, resources *pdf.PdfPageResources, dst, infile string, count *int64) error {
	cstreamParser := pdfcontent.NewContentStreamParser(contents)
	operations, err := cstreamParser.Parse()
	if err != nil {
		return err
	}

	processedXObjects := map[string]bool{}

	// Range through all the content stream operations.
	for _, op := range *operations {
		if op.Operand == "BI" && len(op.Params) == 1 {
			// BI: Inline image.

			iimg, ok := op.Params[0].(*pdfcontent.ContentStreamInlineImage)
			if !ok {
				continue
			}

			img, err := iimg.ToImage(resources)
			if err != nil {
				return err
			}

			cs, err := iimg.GetColorSpace(resources)
			if err != nil {
				return err
			}
			if cs == nil {
				// Default if not specified?
				cs = pdf.NewPdfColorspaceDeviceGray()
			}
			fmt.Printf("Cs: %T\n", cs)

			rgbImg, err := cs.ImageToRGB(*img)
			if err != nil {
				return err
			}

			gimg, err := rgbImg.ToGoImage()
			if err != nil {
				return err
			}

			if err := saveimage(dst, infile, gimg, count); err != nil {
				return err
			}
		} else if op.Operand == "Do" && len(op.Params) == 1 {
			// Do: XObject.
			name := op.Params[0].(*pdfcore.PdfObjectName)

			// Only process each one once.
			_, has := processedXObjects[string(*name)]
			if has {
				continue
			}
			processedXObjects[string(*name)] = true

			_, xtype := resources.GetXObjectByName(*name)
			if xtype == pdf.XObjectTypeImage {
				ximg, err := resources.GetXObjectImageByName(*name)
				if err != nil {
					return err
				}

				img, err := ximg.ToImage()
				if err != nil {
					return err
				}

				rgbImg, err := ximg.ColorSpace.ImageToRGB(*img)
				if err != nil {
					return err
				}

				gimg, err := rgbImg.ToGoImage()
				if err != nil {
					return err
				}

				saveimage(dst, infile, gimg, count)
			} else if xtype == pdf.XObjectTypeForm {
				// Go through the XObject Form content stream.
				xform, err := resources.GetXObjectFormByName(*name)
				if err != nil {
					return err
				}

				formContent, err := xform.GetContentStream()
				if err != nil {
					return err
				}

				// Process the content stream in the Form object too:
				formResources := xform.Resources
				if formResources == nil {
					formResources = resources
				}

				// Process the content stream in the Form object too:
				if err := extImgsInContStream(string(formContent), formResources, dst, infile, count); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func saveimage(dst, infile string, img image.Image, count *int64) error {
	_, fname := filepath.Split(infile)
	name := fmt.Sprintf("%s_%s.jpg", fname, getimgname(4))

	dstFile, err := os.Create(dst + name)
	defer dstFile.Close()
	if err != nil {
		return err
	}

	if err := png.Encode(dstFile, img); err != nil {
		return err
	}

	*count++
	fmt.Printf("\rImage File Found: %s", name)
	return nil
}

func copyimage(buf *[]byte, count *int64, dst, name string) error {
	var err error
	filetypedir(dst)

	if filetype.IsImage(*buf) {
		*count++
		err = ioutil.WriteFile(dst+name, *buf, 0644)
		fmt.Printf("\rImage File Found: %s", name)
	}

	return err
}

func copyvideo(buf *[]byte, count *int64, dst, name string) error {
	var err error
	filetypedir(dst)

	if filetype.IsVideo(*buf) {
		*count++
		err = ioutil.WriteFile(dst+name, *buf, 0644)
		fmt.Printf("\rVideo File Found: %s", name)
	}

	return err
}

func copyaudio(buf *[]byte, count *int64, dst, name string) error {
	var err error
	filetypedir(dst)

	if filetype.IsAudio(*buf) {
		*count++
		err = ioutil.WriteFile(dst+name, *buf, 0644)
		fmt.Printf("\rAudio File Found: %s", name)
	}

	return err
}

func copyarchive(buf *[]byte, count *int64, dst, name string) error {
	var err error
	filetypedir(dst)

	if filetype.IsArchive(*buf) {
		*count++
		err = ioutil.WriteFile(dst+name, *buf, 0644)
		fmt.Printf("\rArchive File Found: %s", name)
	}

	return err
}

func filetypedir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}
}

func ispdf(buf []byte) bool {
	return len(buf) > 3 &&
		buf[0] == 0x25 && buf[1] == 0x50 &&
		buf[2] == 0x44 && buf[3] == 0x46
}

func isepub(buf []byte) bool {
	return len(buf) > 57 &&
		buf[0] == 0x50 && buf[1] == 0x4B && buf[2] == 0x3 && buf[3] == 0x4 &&
		buf[30] == 0x6D && buf[31] == 0x69 && buf[32] == 0x6D && buf[33] == 0x65 &&
		buf[34] == 0x74 && buf[35] == 0x79 && buf[36] == 0x70 && buf[37] == 0x65 &&
		buf[38] == 0x61 && buf[39] == 0x70 && buf[40] == 0x70 && buf[41] == 0x6C &&
		buf[42] == 0x69 && buf[43] == 0x63 && buf[44] == 0x61 && buf[45] == 0x74 &&
		buf[46] == 0x69 && buf[47] == 0x6F && buf[48] == 0x6E && buf[49] == 0x2F &&
		buf[50] == 0x65 && buf[51] == 0x70 && buf[52] == 0x75 && buf[53] == 0x62 &&
		buf[54] == 0x2B && buf[55] == 0x7A && buf[56] == 0x69 && buf[57] == 0x70
}

func getimgname(length int) string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base32.StdEncoding.EncodeToString(randomBytes)[:length]
}

func carvefile(count *int64, dst, path string) error {
	filetypedir(dst)
	_, name := filepath.Split(path)
	fmt.Println("\nFile carver running on: ", name)

	finfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	size := finfo.Size()
	if size >= 2e+9 {
		errors.New("File is too big to be processed")
	}

	fd, err := unix.Open(path, unix.O_RDONLY, 0777)
	defer unix.Close(fd)
	if err != nil {
		return err
	}

	if err := getjpg(count, fd, size, dst); err != nil {
		return err
	}
	if err := getgif(count, fd, size, dst); err != nil {
		return err
	}
	if err := getpng(count, fd, size, dst); err != nil {
		return err
	}

	return nil
}

func getjpg(count *int64, fd int, size int64, dst string) error {
	buff := make([]byte, 1)
	var counter int8
	var carved []byte

	for i := int64(0); i < size; i++ {
		if _, err := unix.Read(fd, buff); err != nil {
			return err
		}

		switch counter {
		case 0:
			if buff[0] == 0xff {
				carved = append(carved, buff[0])
				counter++
			} else {
				carved = nil
				counter = 0
			}
		case 1:
			if buff[0] == 0xd8 {
				carved = append(carved, buff[0])
				counter++
			} else {
				carved = nil
				counter = 0
			}
		case 2:
			if buff[0] == 0xff {
				carved = append(carved, buff[0])
				counter++
			} else {
				carved = nil
				counter = 0
			}
		case 3:
			if buff[0] == 0xff {
				carved = append(carved, buff[0])
				counter++
			} else {
				carved = append(carved, buff[0])
			}
		case 4:
			if buff[0] == 0xd9 {
				carved = append(carved, buff[0])
				if err := writecarved(dst, "jpg", &carved, count); err != nil {
					return err
				}
				carved = nil
				counter = 0
			} else {
				carved = append(carved, buff[0])
				counter--
			}
		}
	}
	return nil
}

func getgif(count *int64, fd int, size int64, dst string) error {
	buff := make([]byte, 1)
	var counter int8
	var carved []byte

	if _, err := unix.Seek(fd, 0, 0); err != nil {
		return err
	}

	for i := int64(0); i < size; i++ {
		if _, err := unix.Read(fd, buff); err != nil {
			return err
		}

		switch counter {
		case 0:
			if buff[0] == 0x47 {
				carved = append(carved, buff[0])
				counter++
			} else {
				carved = nil
				counter = 0
			}
		case 1:
			if buff[0] == 0x49 {
				carved = append(carved, buff[0])
				counter++
			} else {
				carved = nil
				counter = 0
			}
		case 2:
			if buff[0] == 0x46 {
				carved = append(carved, buff[0])
				counter++
			} else {
				carved = nil
				counter = 0
			}
		case 3:
			if buff[0] == 0x00 {
				carved = append(carved, buff[0])
				counter++
			} else {
				carved = append(carved, buff[0])
			}
		case 4:
			if buff[0] == 0x00 {
				carved = append(carved, buff[0])
				counter++
			} else {
				carved = append(carved, buff[0])
				counter--
			}
		case 5:
			if buff[0] == 0x3b {
				carved = append(carved, buff[0])
				if err := writecarved(dst, "gif", &carved, count); err != nil {
					return err
				}
				carved = nil
				counter = 0
			} else {
				carved = append(carved, buff[0])
				counter -= 2
			}
		}
	}
	return nil
}

func getpng(count *int64, fd int, size int64, dst string) error {
	buff := make([]byte, 1)
	var counter int8
	var carved []byte

	if _, err := unix.Seek(fd, 0, 0); err != nil {
		return err
	}

	switch counter {
	case 0:
		if buff[0] == 0x89 {
			carved = append(carved, buff[0])
			counter++
		} else {
			carved = nil
			counter = 0
		}
	case 1:
		if buff[0] == 0x50 {
			carved = append(carved, buff[0])
			counter++
		} else {
			carved = nil
			counter = 0
		}
	case 2:
		if buff[0] == 0x4e {
			carved = append(carved, buff[0])
			counter++
		} else {
			carved = nil
			counter = 0
		}
	case 3:
		if buff[0] == 0x47 {
			carved = append(carved, buff[0])
			counter++
		} else {
			carved = nil
			counter = 0
		}
	case 4:
		if buff[0] == 0xae {
			carved = append(carved, buff[0])
			counter++
		} else {
			carved = append(carved, buff[0])
		}
	case 5:
		if buff[0] == 0x42 {
			carved = append(carved, buff[0])
			counter++
		} else {
			carved = append(carved, buff[0])
			counter--
		}
	case 6:
		if buff[0] == 0x60 {
			carved = append(carved, buff[0])
			counter++
		} else {
			carved = append(carved, buff[0])
			counter -= 2
		}
	case 7:
		if buff[0] == 0x82 {
			carved = append(carved, buff[0])
			if err := writecarved(dst, "png", &carved, count); err != nil {
				return err
			}
			carved = nil
			counter = 0
		} else {
			carved = append(carved, buff[0])
			counter -= 3
		}
	}

	return nil
}

func writecarved(dst, ext string, data *[]byte, count *int64) error {
	name := dst + getimgname(10) + "." + ext
	if err := ioutil.WriteFile(name, *data, 0644); err != nil {
		return err
	}
	*count++
	_, fname := filepath.Split(name)
	fmt.Printf("\rImage file found: %s", fname)
	return nil
}

func getimgfrompop(buf *[]byte, count *int64, dst, path string) error {
	filetypedir(dst)

	if ispopdoc(buf, filepath.Ext(path)) {
		r, err := zip.OpenReader(path)
		defer r.Close()
		if err != nil {
			return err
		}

		for _, f := range r.File {
			if f.FileInfo().IsDir() {
				continue
			} else if isimgext(filepath.Ext(f.Name)) {
				fpath := dst + "/" + filepath.Base(path) + "_" + getimgname(6) + "_" + filepath.Base(f.Name)

				outFile, _ := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
				rc, _ := f.Open()
				io.Copy(outFile, rc)

				outFile.Close()
				rc.Close()
				*count++
				fmt.Printf("\rImage file found: %s", filepath.Base(fpath))
			}
		}
	}

	return nil
}

func isimgext(ext string) bool {
	validexts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".webp": true,
		".svg":  true,
	}
	return validexts[ext]
}

func ispopdoc(buf *[]byte, ext string) bool {
	validext := map[string]bool{
		".odf":   true,
		".odt":   true,
		".odp":   true,
		".pages": true,
	}
	return filetype.IsDocument(*buf) || isepub(*buf) || validext[ext]
}
