package logger

import (
	"fmt"
	"github.com/byawitz/gint/internal/theme"
)

func Notice(messages ...string) {
	fmt.Println(theme.Orange.Render(messages...))
}

func Good(messages ...string) {
	fmt.Println(theme.Green.Render(messages...))
}

func Bad(messages ...string) {
	fmt.Println(theme.Red.Render(messages...))
}

func Fatal(messages ...string) {
	panic(theme.Red.Render(messages...))
}
