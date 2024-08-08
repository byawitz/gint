package commands

import (
	"fmt"
	"github.com/byawitz/gint/internal/configurator"
	"github.com/byawitz/gint/pkg/lexer"
	"os"
)

func Format(ci bool, files []string, config *configurator.Config) {
	ch := make(chan bool)

	for _, file := range files {
		go func() {
			content, err := os.ReadFile(file)
			if err != nil {
				ch <- false
				return
			}
			lexer.Tokenize(string(content))
			ch <- true
			return
		}()
	}

	for _ = range files {
		val := <-ch
		if val {
			fmt.Print("V")
		} else {
			fmt.Print("X")
		}
	}
	fmt.Printf("lexerd through %d files", len(files))
}
