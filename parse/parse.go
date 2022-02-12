package parse

import (
	"errors"
	"strconv"

	"github.com/yuniruyuni/lang/ast"
	"github.com/yuniruyuni/lang/token"
	"github.com/yuniruyuni/lang/token/kind"
)

// --- PEG ---
// Root := Sum | Integer | String
// Sum := Integer `kind.Plus` Sum
// Integer := `kind.Integer`
// String := `kind.String`

type Parser struct {
	cur    int
	tokens []*token.Token
}

func (p *Parser) Next() *token.Token {
	p.cur += 1
	return p.tokens[p.cur-1]
}

func (p *Parser) LookAt(n int) *token.Token {
	if p.cur+n >= len(p.tokens) {
		return nil
	}

	return p.tokens[p.cur+n]
}

func (p *Parser) Root() (ast.AST, error) {
	if len(p.tokens) != 1 {
		return nil, errors.New("not enough tokens")
	}

	t := p.Next()

	cands := []func(t *token.Token) (ast.AST, error){
		p.Sum,
		p.Integer,
		p.String,
	}

	for _, cand := range cands {
		parsed, err := cand(t)
		if err != nil {
			continue
		}
		return parsed, nil
	}

	return nil, errors.New("invalid tokens")
}

func (p *Parser) Integer(t *token.Token) (ast.AST, error) {
	if t.Kind != kind.Integer {
		return nil, errors.New("invalid token")
	}

	val, err := strconv.Atoi(t.Str)
	if err != nil {
		return nil, err
	}
	return &ast.Integer{Value: val}, nil
}

func (p *Parser) String(t *token.Token) (ast.AST, error) {
	if t.Kind != kind.String {
		return nil, errors.New("invalid token")
	}
	return &ast.String{Word: t.Str}, nil
}

func (p *Parser) Sum(t1 *token.Token) (ast.AST, error) {
	lhs, err := p.Integer(t1)
	if err != nil {
		return nil, err
	}

	t2 := p.LookAt(0)
	if t2 == nil || t2.Kind != kind.Plus {
		return lhs, nil
	}

	_ = p.Next()
	t3 := p.Next()

	rhs, err := p.Sum(t3)
	if err != nil {
		return lhs, nil
	}

	res := &ast.Sum{
		LHS: lhs,
		RHS: rhs,
	}

	return res, nil
}

func Parse(tks []*token.Token) (ast.AST, error) {
	parser := Parser{cur: 0, tokens: tks}
	return parser.Root()
}
