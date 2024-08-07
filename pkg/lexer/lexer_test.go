package lexer

import (
	"fmt"
	"os"
	"testing"
)

func TestParsing(t *testing.T) {
	file, err := os.ReadFile("../../tests_assets/lexer/all_tokens.php")
	if err != nil {
		t.Fatal("error opening PHP file for testing")
	}

	tokens := Tokenize(string(file))
	fmt.Println(tokens)
}
