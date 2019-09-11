package lib

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/h2non/filetype"
	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdf "github.com/unidoc/unidoc/pdf/model"
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
				copyimage(&buf, &count, dst+"images/", info.Name())
				getimgfrompdf(&buf, &count, dst+"images/", filepath)
			case VIDEO:
				copyvideo(&buf, &count, dst+"videos/", info.Name())
			case AUDIO:
				copyaudio(&buf, &count, dst+"audios/", info.Name())
			case ARCHIVE:
				copyarchive(&buf, &count, dst+"archives/", info.Name())
			default:
				fmt.Println("Wrong Choice")
				return nil
			}

			return err
		}
		return nil
	})

	if count == 0 {
		fmt.Printf("\nFile not found . - .")
	}
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

	return extractImagesInContentStream(contents, page.Resources, dst, infile, count)
}

func extractImagesInContentStream(contents string, resources *pdf.PdfPageResources, dst, infile string, count *int64) error {
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
				if err := extractImagesInContentStream(string(formContent), formResources, dst, infile, count); err != nil {
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
	fmt.Printf("\rImage File Found: %s, Count: %v", name, *count)
	return nil
}

func copyimage(buf *[]byte, count *int64, dst, name string) error {
	var err error
	filetypedir(dst)

	if filetype.IsImage(*buf) {
		*count++
		err = ioutil.WriteFile(dst+name, *buf, 0644)
		fmt.Printf("\rImage File Found: %s, Count: %v", name, *count)
	}

	return err
}

func copyvideo(buf *[]byte, count *int64, dst, name string) error {
	var err error
	filetypedir(dst)

	if filetype.IsVideo(*buf) {
		*count++
		err = ioutil.WriteFile(dst+name, *buf, 0644)
		fmt.Printf("\rVideo File Found: %s, Count: %v", name, *count)
	}

	return err
}

func copyaudio(buf *[]byte, count *int64, dst, name string) error {
	var err error
	filetypedir(dst)

	if filetype.IsAudio(*buf) {
		*count++
		err = ioutil.WriteFile(dst+name, *buf, 0644)
		fmt.Printf("\rAudio File Found: %s, Count: %v", name, *count)
	}

	return err
}

func copyarchive(buf *[]byte, count *int64, dst, name string) error {
	var err error
	filetypedir(dst)

	if filetype.IsArchive(*buf) {
		*count++
		err = ioutil.WriteFile(dst+name, *buf, 0644)
		fmt.Printf("\rArchive File Found: %s, Count: %v", name, *count)
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

func getimgname(length int) string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base32.StdEncoding.EncodeToString(randomBytes)[:length]
}
