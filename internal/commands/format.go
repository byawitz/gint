package commands

import (
	"fmt"
)

func Format(ci, dirty bool, config string) {
	fmt.Printf("Format %v, %s", []bool{ci, dirty}, config)
}
