package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/h2non/filetype"
)

// Extract function takes in root dir, destination and choice.
// Walks through entire file structure to copy files based on their magic numbers
func Extract(root, dst string, in int) int64 {
	var count int64
	filepath.Walk(root, func(filepath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !info.IsDir() {
			buf, err := ioutil.ReadFile(filepath)

			switch in {
			case IMAGE:
				copyimage(&buf, &count, dst+"images/", info.Name())
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

	return count
}

func copyimage(buf *[]byte, count *int64, dst, name string) error {
	var err error
	filetypedir(dst)

	if filetype.IsImage(*buf) {
		*count++
		err = ioutil.WriteFile(dst+name, *buf, 0644)
		fmt.Printf("\rImage File Found: %s, Count: %v", name, count)
	}

	return err
}

func copyvideo(buf *[]byte, count *int64, dst, name string) error {
	var err error
	filetypedir(dst)

	if filetype.IsVideo(*buf) {
		*count++
		err = ioutil.WriteFile(dst+name, *buf, 0644)
		fmt.Printf("\rVideo File Found: %s, Count: %v", name, count)
	}

	return err
}

func copyaudio(buf *[]byte, count *int64, dst, name string) error {
	var err error
	filetypedir(dst)

	if filetype.IsAudio(*buf) {
		*count++
		err = ioutil.WriteFile(dst+name, *buf, 0644)
		fmt.Printf("\rAudio File Found: %s, Count: %v", name, count)
	}

	return err
}

func copyarchive(buf *[]byte, count *int64, dst, name string) error {
	var err error
	filetypedir(dst)

	if filetype.IsArchive(*buf) {
		*count++
		err = ioutil.WriteFile(dst+name, *buf, 0644)
		fmt.Printf("\rArchive File Found: %s, Count: %v", name, count)
	}

	return err
}

func filetypedir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}
}
