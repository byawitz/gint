package commands

import (
	"fmt"
	"github.com/byawitz/gint/internal/theme"
)

func Version() {
	fmt.Printf("gint %v", theme.Blue.Render("GINT_VERSION"))
}
