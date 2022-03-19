package parse

import (
	"errors"
	"strconv"

	"github.com/yuniruyuni/lang/ast"
	"github.com/yuniruyuni/lang/token"
	"github.com/yuniruyuni/lang/token/kind"
)

// Pos is the position for parsing code.
type Pos int

// Parser transforms this language into AST.
// --- PEG ---
// AST Emit will happen for x in [x].
// Root := Expr | Res | String
// Expr := Add | Sub | Term
// [Add] := Term + Expr
// [Sub] := Term - Expr
// Term := Mul | Div | Res
// [Mul] := Res * Term
// [Div] := Res / Term
// Res := Clause | Integer
// Clause := ( Expr )
type Parser struct {
	tokens []*token.Token
}

func (p *Parser) Len() Pos {
	return Pos(len(p.tokens))
}

func (p *Parser) LookAt(at Pos) *token.Token {
	if at < p.Len() {
		return p.tokens[at]
	}
	return nil
}

func (p *Parser) Consume(kind kind.Kind, at Pos) (Pos, *token.Token) {
	t := p.LookAt(at)
	if t == nil || t.Kind != kind {
		return at, nil
	}
	return at + 1, t
}

func (p *Parser) End(at Pos) bool {
	return at == p.Len()
}

func (p *Parser) Root(at Pos) (Pos, ast.AST, error) {
	cands := []func(Pos) (Pos, ast.AST, error){
		p.Expr,
		p.Res,
		p.String,
	}

	for _, cand := range cands {
		nx, parsed, err := cand(at)
		if err != nil {
			continue
		}
		if p.End(nx) {
			return nx, parsed, nil
		}
	}

	return at, nil, errors.New("invalid tokens")
}

func (p *Parser) Expr(at Pos) (Pos, ast.AST, error) {
	cands := []func(Pos) (Pos, ast.AST, error){
		p.Add,
		p.Sub,
		p.Term,
	}

	for _, cand := range cands {
		nx, parsed, err := cand(at)
		if err == nil {
			return nx, parsed, nil
		}
	}

	return at, nil, errors.New("invalid token")
}

func (p *Parser) Add(at Pos) (Pos, ast.AST, error) {
	nx := at
	nx, lhs, err := p.Term(nx)
	if err != nil {
		return at, nil, err
	}

	nx, pls := p.Consume(kind.Plus, nx)
	if pls == nil {
		return at, nil, errors.New("invalid tokens")
	}

	nx, rhs, err := p.Expr(nx)
	if err != nil {
		return at, nil, err
	}

	return nx, &ast.Add{LHS: lhs, RHS: rhs}, nil
}

func (p *Parser) Sub(at Pos) (Pos, ast.AST, error) {
	nx := at
	nx, lhs, err := p.Term(nx)
	if err != nil {
		return at, nil, err
	}

	nx, mns := p.Consume(kind.Minus, nx)
	if mns == nil {
		return at, nil, errors.New("invalid tokens")
	}

	nx, rhs, err := p.Expr(nx)
	if err != nil {
		return at, nil, err
	}

	return nx, &ast.Sub{LHS: lhs, RHS: rhs}, nil
}

func (p *Parser) Term(at Pos) (Pos, ast.AST, error) {
	cands := []func(Pos) (Pos, ast.AST, error){
		p.Mul,
		p.Div,
		p.Res,
	}

	for _, cand := range cands {
		nx, parsed, err := cand(at)
		if err == nil {
			return nx, parsed, nil
		}
	}

	return at, nil, errors.New("invalid token")
}

func (p *Parser) Mul(at Pos) (Pos, ast.AST, error) {
	nx := at
	nx, lhs, err := p.Res(nx)
	if err != nil {
		return at, nil, err
	}

	nx, mul := p.Consume(kind.Multiply, nx)
	if mul == nil {
		return at, nil, errors.New("invalid tokens")
	}

	nx, rhs, err := p.Term(nx)
	if err != nil {
		return at, nil, err
	}

	return nx, &ast.Mul{LHS: lhs, RHS: rhs}, nil
}

func (p *Parser) Div(at Pos) (Pos, ast.AST, error) {
	nx := at
	nx, lhs, err := p.Res(nx)
	if err != nil {
		return at, nil, err
	}

	nx, div := p.Consume(kind.Divide, nx)
	if div == nil {
		return at, nil, errors.New("invalid tokens")
	}

	nx, rhs, err := p.Term(nx)
	if err != nil {
		return at, nil, err
	}

	return nx, &ast.Div{LHS: lhs, RHS: rhs}, nil
}

func (p *Parser) Res(at Pos) (Pos, ast.AST, error) {
	cands := []func(Pos) (Pos, ast.AST, error){
		p.Clause,
		p.Integer,
	}

	for _, cand := range cands {
		nx, parsed, err := cand(at)
		if err == nil {
			return nx, parsed, nil
		}
	}

	return at, nil, errors.New("invalid token")
}

func (p *Parser) Clause(at Pos) (Pos, ast.AST, error) {
	nx := at
	nx, lp := p.Consume(kind.LeftParen, nx)
	if lp == nil {
		return at, nil, errors.New("invalid token")
	}

	nx, child, err := p.Expr(nx)
	if err != nil {
		return at, nil, err
	}

	nx, rp := p.Consume(kind.RightParen, nx)
	if rp == nil {
		return at, nil, errors.New("invalid token")
	}

	return nx, child, nil
}

func (p *Parser) Integer(at Pos) (Pos, ast.AST, error) {
	nx := at
	nx, t := p.Consume(kind.Integer, nx)
	if t == nil {
		return at, nil, errors.New("invalid token")
	}

	val, err := strconv.Atoi(t.Str)
	if err != nil {
		return at, nil, errors.New("Integer constant size over than max bit size")
	}
	return nx, &ast.Integer{Value: val}, nil
}

func (p *Parser) String(at Pos) (Pos, ast.AST, error) {
	nx, t := p.Consume(kind.String, at)
	if t == nil {
		return at, nil, errors.New("invalid token")
	}
	return nx, &ast.String{Word: t.Str}, nil
}

func Parse(tks []*token.Token) (ast.AST, error) {
	parser := Parser{tokens: tks}
	_, ast, err := parser.Root(0)
	return ast, err
}
