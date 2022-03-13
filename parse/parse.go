package parse

import (
	"errors"
	"strconv"

	"github.com/yuniruyuni/lang/ast"
	"github.com/yuniruyuni/lang/token"
	"github.com/yuniruyuni/lang/token/kind"
)

// Parser transforms this language into AST.
// --- PEG ---
// AST Emit will happen for x in [x].
// Root := Expr | Res | String
// Expr := Term Add | Term Sub | Term
// [Add] := + Expr
// [Sub] := - Expr
// Term := Res Mul | Res Div | Res
// [Mul] := * Term
// [Div] := / Term
// Res := Clause | Integer
// Clause := ( Expr )
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
		p.Expr,
		p.Res,
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

func (p *Parser) Expr() (ast.AST, error) {
	lhs, err := p.Term()
	if err != nil {
		return nil, err
	}

	cands := []func(ast.AST) (ast.AST, error){
		p.Add,
		p.Sub,
	}

	for _, cand := range cands {
		parsed, err := cand(lhs)
		if err == nil {
			return parsed, nil
		}
	}

	return lhs, nil
}

func (p *Parser) Add(lhs ast.AST) (ast.AST, error) {
	pls := p.Consume(kind.Plus)
	if pls == nil {
		return nil, errors.New("invalid tokens")
	}

	rhs, err := p.Expr()
	if err != nil {
		p.Revert(1)
		return lhs, nil
	}

	return &ast.Add{
		LHS: lhs,
		RHS: rhs,
	}, nil
}

func (p *Parser) Sub(lhs ast.AST) (ast.AST, error) {
	mns := p.Consume(kind.Minus)
	if mns == nil {
		return nil, errors.New("invalid tokens")
	}

	rhs, err := p.Expr()
	if err != nil {
		p.Revert(1)
		return lhs, nil
	}

	return &ast.Sub{
		LHS: lhs,
		RHS: rhs,
	}, nil
}

func (p *Parser) Term() (ast.AST, error) {
	lhs, err := p.Res()
	if err != nil {
		return nil, err
	}

	cands := []func(ast.AST) (ast.AST, error){
		p.Mul,
		p.Div,
	}

	for _, cand := range cands {
		parsed, err := cand(lhs)
		if err == nil {
			return parsed, nil
		}
	}

	return lhs, nil
}

func (p *Parser) Mul(lhs ast.AST) (ast.AST, error) {
	mul := p.Consume(kind.Multiply)
	if mul == nil {
		return nil, errors.New("invalid tokens")
	}

	rhs, err := p.Term()
	if err != nil {
		p.Revert(1)
		return lhs, nil
	}

	return &ast.Mul{
		LHS: lhs,
		RHS: rhs,
	}, nil
}

func (p *Parser) Div(lhs ast.AST) (ast.AST, error) {
	div := p.Consume(kind.Divide)
	if div == nil {
		return nil, errors.New("invalid tokens")
	}

	rhs, err := p.Term()
	if err != nil {
		p.Revert(1)
		return lhs, nil
	}

	return &ast.Div{
		LHS: lhs,
		RHS: rhs,
	}, nil
}

func (p *Parser) Res() (ast.AST, error) {
	cands := []func() (ast.AST, error){
		p.Clause,
		p.Integer,
	}

	for _, cand := range cands {
		parsed, err := cand()
		if err == nil {
			return parsed, nil
		}
	}

	return nil, errors.New("invalid token")
}

func (p *Parser) Clause() (ast.AST, error) {
	lp := p.Consume(kind.LeftParen)
	if lp == nil {
		return nil, errors.New("invalid token")
	}

	child, err := p.Expr()
	if err != nil {
		p.Revert(1)
		return nil, err
	}

	rp := p.Consume(kind.RightParen)
	if rp == nil {
		p.Revert(2)
		return nil, errors.New("invalid token")
	}

	return child, nil
}

func (p *Parser) Integer() (ast.AST, error) {
	t := p.Consume(kind.Integer)
	if t == nil {
		return nil, errors.New("invalid token")
	}

	val, err := strconv.Atoi(t.Str)
	if err != nil {
		return nil, errors.New("Integer constant size over than max bit size")
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

func Parse(tks []*token.Token) (ast.AST, error) {
	parser := Parser{cur: 0, tokens: tks}
	return parser.Root()
}
