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
				if filetype.IsImage(buf) {
					count++
					if err := ioutil.WriteFile(dst+info.Name(), buf, 0644); err != nil {
						fmt.Println(err)
						return err
					}
					fmt.Printf("\rImage File Found: %s, Count: %d", info.Name(), count)
				}
			case VIDEO:
				if filetype.IsVideo(buf) {
					count++
					if err := ioutil.WriteFile(dst+info.Name(), buf, 0644); err != nil {
						fmt.Println(err)
						return err
					}
					fmt.Printf("\rVideo File Found: %s, Count: %d", info.Name(), count)
				}
			case AUDIO:
				if filetype.IsAudio(buf) {
					count++
					if err := ioutil.WriteFile(dst+info.Name(), buf, 0644); err != nil {
						fmt.Println(err)
						return err
					}
					fmt.Printf("\rAudio File Found: %s, Count: %d", info.Name(), count)
				}
			case ARCHIVE:
				if filetype.IsArchive(buf) {
					count++
					if err := ioutil.WriteFile(dst+info.Name(), buf, 0644); err != nil {
						fmt.Println(err)
						return err
					}
					fmt.Printf("\rArchive File Found: %s, Count: %d", info.Name(), count)
				}
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
