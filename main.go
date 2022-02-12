package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/yuniruyuni/lang/ast"
	"github.com/yuniruyuni/lang/gen"
	"github.com/yuniruyuni/lang/token"
	"github.com/yuniruyuni/lang/token/kind"
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

func parse(tks []*token.Token) (ast.AST, error) {
	if len(tks) != 1 {
		return nil, errors.New("not enough tokens")
	}

	t := tks[0]
	switch t.Kind {
	case kind.String:
		return &ast.String{Word: t.Str}, nil
	case kind.Integer:
		val, err := strconv.Atoi(t.Str)
		if err != nil {
			return nil, err
		}
		return &ast.Integer{Value: val}, nil
	default:
		return nil, errors.New("unsupported node")
	}
}

func main() {
	sc.Scan()
	code := sc.Text()

	tks, err := tokenize(code)
	if err != nil {
		_ = fmt.Errorf("failed to tokenize code: %s", err.Error())
		os.Exit(-1)
	}

	root, err := parse(tks)
	if err != nil {
		_ = fmt.Errorf("failed to parse code: %s", err.Error())
		os.Exit(-1)
	}

	outputLL(root)
}
