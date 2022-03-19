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

// NonTerminal expresses non-terminal symbol in parser.
type NonTerminal func(Pos) (Pos, ast.AST, error)

// Select combines some target NonTerminals into single NonTerminal.
// This new NonTerminal checks if targets match current tokens and
// returns first maching NonTerminal result.
// If there is no maching NonTerminal, It returns an "invalid tokens" error.
func Select(cands ...NonTerminal) NonTerminal {
	return func(at Pos) (Pos, ast.AST, error) {
		for _, cand := range cands {
			nx, parsed, err := cand(at)
			if err == nil {
				return nx, parsed, nil
			}
		}
		return at, nil, errors.New("invalid tokens")
	}
}

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

func (p *Parser) Skip(kind kind.Kind) NonTerminal {
	return func(at Pos) (Pos, ast.AST, error) {
		nx, t := p.Consume(kind, at)
		if t == nil {
			return at, nil, errors.New("invalid token")
		}
		return nx, nil, nil
	}
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
	nx, parsed, err := Select(p.Expr, p.Res, p.String)(at)
	if err != nil {
		return at, nil, err
	}

	if !p.End(nx) {
		return at, nil, errors.New("invalid tokens")
	}

	return nx, parsed, nil
}

func (p *Parser) Expr(at Pos) (Pos, ast.AST, error) {
	return Select(p.Add, p.Sub, p.Term)(at)
}

func (p *Parser) Add(at Pos) (Pos, ast.AST, error) {
	nx := at
	nx, lhs, err := p.Term(nx)
	if err != nil {
		return at, nil, err
	}

	nx, _, err = p.Skip(kind.Plus)(nx)
	if err != nil {
		return at, nil, err
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

	nx, _, err = p.Skip(kind.Minus)(nx)
	if err != nil {
		return at, nil, err
	}

	nx, rhs, err := p.Expr(nx)
	if err != nil {
		return at, nil, err
	}

	return nx, &ast.Sub{LHS: lhs, RHS: rhs}, nil
}

func (p *Parser) Term(at Pos) (Pos, ast.AST, error) {
	return Select(p.Mul, p.Div, p.Res)(at)
}

func (p *Parser) Mul(at Pos) (Pos, ast.AST, error) {
	nx := at
	nx, lhs, err := p.Res(nx)
	if err != nil {
		return at, nil, err
	}

	nx, _, err = p.Skip(kind.Multiply)(nx)
	if err != nil {
		return at, nil, err
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

	nx, _, err = p.Skip(kind.Divide)(nx)
	if err != nil {
		return at, nil, err
	}

	nx, rhs, err := p.Term(nx)
	if err != nil {
		return at, nil, err
	}

	return nx, &ast.Div{LHS: lhs, RHS: rhs}, nil
}

func (p *Parser) Res(at Pos) (Pos, ast.AST, error) {
	return Select(p.Clause, p.Integer)(at)
}

func (p *Parser) Clause(at Pos) (Pos, ast.AST, error) {
	nx := at

	nx, _, err := p.Skip(kind.LeftParen)(nx)
	if err != nil {
		return at, nil, err
	}

	nx, child, err := p.Expr(nx)
	if err != nil {
		return at, nil, err
	}

	nx, _, err = p.Skip(kind.RightParen)(nx)
	if err != nil {
		return at, nil, err
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
