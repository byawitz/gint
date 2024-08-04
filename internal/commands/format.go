package commands

import (
	"fmt"
)

func Format(ci bool, files []string, config string) {
	fmt.Printf("Format %v, %v, %s", ci, "files", config)
}
