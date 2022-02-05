package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/yuniruyuni/lang/gen"
	"github.com/yuniruyuni/lang/token"
)

func specifiedString(code string) (string, error) {
	t := token.Tokenizer{}
	tks := t.Tokenize(code)
	if len(tks) == 0 {
		return "", errors.New("failed to tokenize")
	}
	return tks[0].Str, nil
}

func outputLL(word string) {
	ll := gen.LLFile{Word: word}
	fmt.Print(ll.Generate())
}

var sc = bufio.NewScanner(os.Stdin)

func main() {
	sc.Scan()
	code := sc.Text()
	word, err := specifiedString(code)
	if err != nil {
		_ = fmt.Errorf("failed to tokenize code.")
		os.Exit(-1)
	}
	outputLL(word)
}
