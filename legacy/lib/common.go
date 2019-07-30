package lib

import "fmt"

// Handle is the common error handling fucton
func Handle(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
