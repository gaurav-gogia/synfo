package lib

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"path/filepath"

	"github.com/h2non/filetype"
	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func Extract(root, dst string, ft string) {
	var count int64
	filepath.Walk(root, func(filepath string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Print("\n", err, "\n\n")
			return nil
		}

		if !info.IsDir() {
			buf, err := os.ReadFile(filepath)
			switch ft {
			case IMAGE:
				if err := copyimage(dst+"images/", info.Name(), &buf, &count); err != nil {
					fmt.Println(err)
				}
				if err := getimgfrompdf(dst+"images/", filepath, &buf, &count); err != nil {
					fmt.Println(err)
				}
				if err := carvefile(dst+"images/", filepath, &buf, &count); err != nil {
					fmt.Println(err)
				}
				if err := getimgfrompop(dst+"images/", filepath, &buf, &count); err != nil {
					fmt.Println(err)
				}
			case VIDEO:
				if err := copyvideo(dst+"videos/", info.Name(), &buf, &count); err != nil {
					fmt.Println(err)
				}
			case AUDIO:
				if err := copyaudio(dst+"audios/", info.Name(), &buf, &count); err != nil {
					fmt.Println(err)
				}
			case ARCHIVE:
				if err := copyarchive(dst+"archives/", info.Name(), &buf, &count); err != nil {
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

	fmt.Println("\n\nTotal files found: ", count)
}

func getimgfrompdf(dstPath, srcFile string, buf *[]byte, count *int64) error {
	filetypedir(dstPath)

	if ispdf(*buf) {
		f, err := os.Open(srcFile)
		if err != nil {
			return err
		}
		defer f.Close()

		pdfreader, err := pdf.NewPdfReader(f)
		if err != nil {
			return err
		}

		isenc, err := pdfreader.IsEncrypted()
		if err != nil {
			return err
		}

		if isenc {
			return errors.New("file is encrypted")
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

			if err := extractImagesOnPage(page, dstPath, srcFile, count); err != nil {
				return err
			}
		}
	}

	return nil
}

func extractImagesOnPage(page *pdf.PdfPage, dstPath, srcFile string, count *int64) error {
	contents, err := page.GetAllContentStreams()
	if err != nil {
		return err
	}

	return extImgsInContStream(contents, page.Resources, dstPath, srcFile, count)
}

func extImgsInContStream(contents string, resources *pdf.PdfPageResources, dstPath, srcFile string, count *int64) error {
	cstreamParser := pdfcontent.NewContentStreamParser(contents)
	operations, err := cstreamParser.Parse()
	if err != nil {
		return err
	}

	processedXObjects := map[string]bool{}

	for _, op := range *operations {
		if op.Operand == "BI" && len(op.Params) == 1 {
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

			if err := saveimage(dstPath, srcFile, gimg, count); err != nil {
				return err
			}
		} else if op.Operand == "Do" && len(op.Params) == 1 {
			name := op.Params[0].(*pdfcore.PdfObjectName)

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

				saveimage(dstPath, srcFile, gimg, count)
			} else if xtype == pdf.XObjectTypeForm {
				xform, err := resources.GetXObjectFormByName(*name)
				if err != nil {
					return err
				}

				formContent, err := xform.GetContentStream()
				if err != nil {
					return err
				}

				formResources := xform.Resources
				if formResources == nil {
					formResources = resources
				}

				if err := extImgsInContStream(string(formContent), formResources, dstPath, srcFile, count); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func saveimage(dstPath, infile string, img image.Image, count *int64) error {
	_, fname := filepath.Split(infile)
	name := fmt.Sprintf("%s_%s.jpg", fname, getimgname(4))

	dstFile, err := os.Create(dstPath + name)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if err := png.Encode(dstFile, img); err != nil {
		return err
	}

	*count++
	fmt.Printf("\rImage File Found: %s", name)
	return nil
}

func copyimage(dstPath, name string, buf *[]byte, count *int64) error {
	var err error
	filetypedir(dstPath)

	if filetype.IsImage(*buf) {
		*count++
		err = os.WriteFile(dstPath+name, *buf, 0644)
		fmt.Printf("\rImage File Found: %s", name)
	}

	return err
}

func copyvideo(dstPath, name string, buf *[]byte, count *int64) error {
	var err error
	filetypedir(dstPath)

	if filetype.IsVideo(*buf) {
		*count++
		err = os.WriteFile(dstPath+name, *buf, 0644)
		fmt.Printf("\rVideo File Found: %s", name)
	}

	return err
}

func copyaudio(dstPath, name string, buf *[]byte, count *int64) error {
	var err error
	filetypedir(dstPath)

	if filetype.IsAudio(*buf) {
		*count++
		err = os.WriteFile(dstPath+name, *buf, 0644)
		fmt.Printf("\rAudio File Found: %s", name)
	}

	return err
}

func copyarchive(dstPath, name string, buf *[]byte, count *int64) error {
	var err error
	filetypedir(dstPath)

	if filetype.IsArchive(*buf) {
		*count++
		err = os.WriteFile(dstPath+name, *buf, 0644)
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

func carvefile(dstPath, srcFile string, buf *[]byte, count *int64) error {
	if ispopdoc(buf, filepath.Ext(srcFile)) || ispdf(*buf) || filetype.IsImage(*buf) {
		return nil
	}

	filetypedir(dstPath)
	_, name := filepath.Split(srcFile)
	fmt.Printf("\rFile carver running on: %s", name)

	finfo, err := os.Stat(srcFile)
	if err != nil {
		return err
	}

	size := finfo.Size()
	if size >= 2e+9 {
		return errors.New("file is too big to be processed")
	}

	data, err := os.ReadFile(srcFile)
	if err != nil {
		return err
	}

	if err := getpics(dstPath, "jpg", data, JPG_HEADER, JPG_FOOTER, count); err != nil {
		return err
	}
	if err := getpics(dstPath, "png", data, PNG_HEADER, PNG_FOOTER, count); err != nil {
		return err
	}
	if err := getpics(dstPath, "gif", data, GIF87A_HEADER, GIF_FOOTER, count); err != nil {
		return err
	}
	if err := getpics(dstPath, "gif", data, GIF98A_HEADER, GIF_FOOTER, count); err != nil {
		return err
	}

	return nil
}

func getpics(dstPath, ext string, data, hdr, ftr []byte, count *int64) error {
	start := bytes.Index(data, hdr)
	end := bytes.Index(data, ftr)

	for start != -1 || end != -1 {
		end = end + len(ftr)
		if start < end {
			*count++
			carved := data[start : end-1]
			err := writecarved(dstPath, ext, carved, count)
			if err != nil {
				return err
			}
		}

		data = data[end:]
		start = bytes.Index(data, hdr)
		end = bytes.Index(data, ftr)
	}

	return nil
}

func writecarved(dstPath, ext string, data []byte, count *int64) error {
	name := dstPath + getimgname(10) + "." + ext
	if err := os.WriteFile(name, data, 0644); err != nil {
		return err
	}
	*count++
	_, fname := filepath.Split(name)
	fmt.Printf("\rImage file found: %s", fname)
	return nil
}

func getimgfrompop(dstPath, dstFile string, buf *[]byte, count *int64) error {
	filetypedir(dstPath)

	if ispopdoc(buf, filepath.Ext(dstFile)) {
		r, err := zip.OpenReader(dstFile)
		if err != nil {
			return err
		}
		defer r.Close()

		for _, f := range r.File {
			if f.FileInfo().IsDir() {
				continue
			} else if isimgext(filepath.Ext(f.Name)) {
				fpath := dstPath + "/" + filepath.Base(dstFile) + "_" + getimgname(6) + "_" + filepath.Base(f.Name)

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
