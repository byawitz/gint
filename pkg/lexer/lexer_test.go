package lexer

import (
	"os"
	"testing"
)

func TestParsing(t *testing.T) {
	file, err := os.ReadFile("../../tests_assets/lexer/all_tokens.php")
	if err != nil {
		t.Fatal("error opening PHP file for testing")
	}

	tokens := Tokenize(string(file))
	for _, token := range tokens {

		token.Debug()
	}
}
