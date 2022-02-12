package parse

import (
	"errors"
	"strconv"

	"github.com/yuniruyuni/lang/ast"
	"github.com/yuniruyuni/lang/token"
	"github.com/yuniruyuni/lang/token/kind"
)

// BNF
// Root := Integer | String | Sum
// Integer := `kind.Integer`
// String := `kind.String`
// Sum := Integer `kind.Plus` Sum
//      | Integer

func Parse(tks []*token.Token) (ast.AST, error) {
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
