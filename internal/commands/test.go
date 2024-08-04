package commands

import (
	"fmt"
	"github.com/byawitz/gint/internal/configurator"
)

func Test(ci bool, files []string, config *configurator.Config) {
	fmt.Printf("Test %v, %v, %s", ci, "files", config)
}
