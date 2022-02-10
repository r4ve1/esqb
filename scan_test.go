package esqb

import (
	"fmt"
	"testing"
)

func Test_parsing(t *testing.T) {
	expr := `ip=="1.1.1.1"`
	tokens, err := scanTokens(expr)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tokens)
}
