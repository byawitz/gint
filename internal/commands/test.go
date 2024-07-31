package commands

import (
	"fmt"
)

func Test(ci, dirty bool, config string) {
	fmt.Printf("Test %v, %s", []bool{ci, dirty}, config)
}
