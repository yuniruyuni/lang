package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/yuniruyuni/lang/ast"
	"github.com/yuniruyuni/lang/gen"
	"github.com/yuniruyuni/lang/parse"
	"github.com/yuniruyuni/lang/token"
)

func outputLL(root ast.AST) string {
	ll := gen.LLFile{AST: root}
	return string(ll.Generate())
}

func tokenize(code string) ([]*token.Token, error) {
	t := token.Tokenizer{}
	tks := t.Tokenize(code)
	if len(tks) == 0 {
		return nil, errors.New("failed to tokenize")
	}
	return tks, nil
}

func Compile(code string) (string, error) {
	tks, err := tokenize(code)
	if err != nil {
		return "", fmt.Errorf("failed to tokenize code: %s", err.Error())
	}

	root, err := parse.Parse(tks)
	if err != nil {
		return "", fmt.Errorf("failed to parse code: %s", err.Error())
	}

	return outputLL(root), nil
}

func main() {
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic("cannot read stdin")
	}
	code := string(bytes)

	ll, err := Compile(code)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(-1)
	}
	fmt.Println(ll)
}
