package commands

import (
	"fmt"
)

func Bail(ci bool, files []string, config string) {
	fmt.Printf("Bail %v, %v, %s", ci, "", config)
}
