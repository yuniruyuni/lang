package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/yuniruyuni/lang/gen"
	"github.com/yuniruyuni/lang/token"
)

func tokenize(s string) []*token.Token {
	t := token.Tokenizer{}
	return t.Tokenize(s)
}

func outputLL(word string) {
	ll := gen.LLFile{Word: word}
	fmt.Print(ll.Generate())
}

var sc = bufio.NewScanner(os.Stdin)

func main() {
	sc.Scan()
	word := sc.Text()
	outputLL(word)
}
