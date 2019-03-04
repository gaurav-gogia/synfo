package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	filetype "gopkg.in/h2non/filetype.v1"
)

// Extract function takes root dir, out dir and selection var
func Extract(root, dst string, in int8) {
	var count int64

	filepath.Walk(root, func(filepath string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("\n", err)
			return nil
		}
		if !info.IsDir() {
			buf, err := ioutil.ReadFile(filepath)

			switch in {
			case 1:
				if filetype.IsImage(buf) {
					count++
					if err := ioutil.WriteFile(dst+info.Name(), buf, 0644); err != nil {
						fmt.Println(err)
						return err
					}
					fmt.Printf("\rImage File Found: %s, Count: %v", info.Name(), count)
				}
			case 2:
				if filetype.IsVideo(buf) {
					count++
					if err := ioutil.WriteFile(dst+info.Name(), buf, 0644); err != nil {
						fmt.Println(err)
						return err
					}
					fmt.Printf("\rVideo File Found: %s, Count: %v", info.Name(), count)
				}
			case 3:
				if filetype.IsAudio(buf) {
					count++
					if err := ioutil.WriteFile(dst+info.Name(), buf, 0644); err != nil {
						fmt.Println(err)
						return err
					}
					fmt.Printf("\rAudio File Found: %s, Count: %v", info.Name(), count)
				}
			case 4:
				if filetype.IsArchive(buf) {
					count++
					if err := ioutil.WriteFile(dst+info.Name(), buf, 0644); err != nil {
						fmt.Println(err)
						return err
					}
					fmt.Printf("\rArchive File Found: %s, Count: %v", info.Name(), count)
				}
			default:
				fmt.Println("Wrong Choice")
				return nil
			}

			return err
		}
		return nil
	})
}
