package commands

import (
	"fmt"
)

func Test(ci bool, files []string, config string) {
	fmt.Printf("Test %v, %v, %s", ci, "files", config)
}
