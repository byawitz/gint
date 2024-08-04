package commands

import (
	"fmt"
	"github.com/byawitz/gint/internal/configurator"
)

func Bail(ci bool, files []string, config *configurator.Config) {
	fmt.Printf("Bail %v, %v, %s", ci, "", config)
}
