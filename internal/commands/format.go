package commands

import (
	"fmt"
	"github.com/byawitz/gint/internal/configurator"
)

func Format(ci bool, files []string, config *configurator.Config) {
	fmt.Printf("Format %v, \n%v\n, %s", ci, files, config)
}
