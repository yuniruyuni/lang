package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/yuniruyuni/lang/ast"
	"github.com/yuniruyuni/lang/gen"
	"github.com/yuniruyuni/lang/parse"
	"github.com/yuniruyuni/lang/token"
)

func outputLL(root ast.AST) {
	ll := gen.LLFile{AST: root}
	fmt.Print(ll.Generate())
}

var sc = bufio.NewScanner(os.Stdin)

func tokenize(code string) ([]*token.Token, error) {
	t := token.Tokenizer{}
	tks := t.Tokenize(code)
	if len(tks) == 0 {
		return nil, errors.New("failed to tokenize")
	}
	return tks, nil
}

func main() {
	sc.Scan()
	code := sc.Text()

	tks, err := tokenize(code)
	if err != nil {
		_ = fmt.Errorf("failed to tokenize code: %s", err.Error())
		os.Exit(-1)
	}

	root, err := parse.Parse(tks)
	if err != nil {
		_ = fmt.Errorf("failed to parse code: %s", err.Error())
		os.Exit(-1)
	}

	outputLL(root)
}
