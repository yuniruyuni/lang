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

func (p *Parser) Advance(n int) {
	p.cur += n
}

func (p *Parser) Revert(n int) {
	p.cur -= n
}

func (p *Parser) LookAt(n int) *token.Token {
	at := p.cur + n
	if at < len(p.tokens) {
		return p.tokens[at]
	}
	return nil
}

func (p *Parser) Consume(kind kind.Kind) *token.Token {
	t := p.LookAt(0)
	if t == nil || t.Kind != kind {
		return nil
	}
	p.Advance(1)
	return t
}

func (p *Parser) End() bool {
	return p.cur == len(p.tokens)
}

func (p *Parser) Root() (ast.AST, error) {
	cands := []func() (ast.AST, error){
		p.Sum,
		p.Integer,
		p.String,
	}

	for _, cand := range cands {
		parsed, err := cand()
		if err != nil {
			continue
		}
		if p.End() {
			return parsed, nil
		}
	}

	return nil, errors.New("invalid tokens")
}

func (p *Parser) Integer() (ast.AST, error) {
	t := p.Consume(kind.Integer)
	if t == nil {
		return nil, errors.New("invalid token")
	}

	val, err := strconv.Atoi(t.Str)
	if err != nil {
		panic("Integer constant size over than max bit size")
	}
	return &ast.Integer{Value: val}, nil
}

func (p *Parser) String() (ast.AST, error) {
	t := p.Consume(kind.String)
	if t == nil {
		return nil, errors.New("invalid token")
	}
	return &ast.String{Word: t.Str}, nil
}

func (p *Parser) Sum() (ast.AST, error) {
	// Check does we have a sequence <Integer, Plus, Sum>

	lhs, err := p.Integer()
	if err != nil {
		return nil, err
	}

	t2 := p.Consume(kind.Plus)
	if t2 == nil {
		return lhs, nil
	}

	rhs, err := p.Sum()
	if err != nil {
		p.Revert(1)
		return lhs, nil
	}

	return &ast.Sum{
		LHS: lhs,
		RHS: rhs,
	}, nil
}

func Parse(tks []*token.Token) (ast.AST, error) {
	parser := Parser{cur: 0, tokens: tks}
	return parser.Root()
}
