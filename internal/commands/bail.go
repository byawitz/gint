package commands

import (
	"fmt"
)

func Bail(ci, dirty bool, config string) {
	fmt.Printf("Bail %v, %s", []bool{ci, dirty}, config)
}
