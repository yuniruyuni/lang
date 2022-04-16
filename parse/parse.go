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

type Key struct {
	Ptr uintptr
	At  Pos
}

type Result struct {
	Pos Pos
	Ast ast.AST
}

type Cache map[Key]*Result

// Parser transforms this language into AST.
// --- PEG ---
// AST Emit will happen for x in [x].
// Root := Exec
// Execute := Sequence | Statement
// [Sequence] := Statement ; Execute
// Statement := Let | Assign | Cond | Res | String
// [Let] := let Variable = Cond
// [Assign] := Variable = Cond
// Cond := Less | Equal | Expr
// [Less] := Expr < Cond
// [Equal] := Expr == Cond
// Expr := Add | Sub | Term
// [Add] := Term + Expr
// [Sub] := Term - Expr
// Term := Mul | Div | Res
// [Mul] := Res * Term
// [Div] := Res / Term
// Res := If | Clause | Variable | Integer
// [Variable] := Identifier
// Clause := ( Cond )
// If := if Execute { Execute } else { Execute }
type Parser struct {
	tokens []*token.Token
	cache  Cache
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
	nx, parsed, err := p.Execute(at)
	if err != nil {
		return at, nil, err
	}

	if !p.End(nx) {
		return at, nil, errors.New("invalid tokens")
	}

	return nx, parsed, nil
}

func (p *Parser) Execute(at Pos) (Pos, ast.AST, error) {
	return p.Select(p.Sequence, p.Statement)(at)
}

func (p *Parser) Sequence(at Pos) (Pos, ast.AST, error) {
	return p.Concat(
		func(asts []ast.AST) ast.AST {
			return &ast.Sequence{LHS: asts[0], RHS: asts[2]}
		},
		p.Statement,
		p.Skip(kind.Semicolon),
		p.Execute,
	)(at)
}

func (p *Parser) Statement(at Pos) (Pos, ast.AST, error) {
	return p.Select(p.Let, p.Assign, p.Cond, p.Res, p.String)(at)
}

func (p *Parser) Let(at Pos) (Pos, ast.AST, error) {
	return p.Concat(
		func(asts []ast.AST) ast.AST {
			return &ast.Let{LHS: asts[1], RHS: asts[3]}
		},
		p.Skip(kind.Let),
		p.Variable,
		p.Skip(kind.Equal),
		p.Cond,
	)(at)
}

func (p *Parser) Assign(at Pos) (Pos, ast.AST, error) {
	return p.Concat(
		func(asts []ast.AST) ast.AST {
			return &ast.Assign{LHS: asts[0], RHS: asts[2]}
		},
		p.Variable,
		p.Skip(kind.Equal),
		p.Cond,
	)(at)
}

func (p *Parser) Cond(at Pos) (Pos, ast.AST, error) {
	return p.Select(p.Less, p.Equal, p.Expr)(at)
}

func (p *Parser) Less(at Pos) (Pos, ast.AST, error) {
	return p.Concat(
		func(asts []ast.AST) ast.AST {
			return &ast.Less{LHS: asts[0], RHS: asts[2]}
		},
		p.Expr,
		p.Skip(kind.Less),
		p.Cond,
	)(at)
}

func (p *Parser) Equal(at Pos) (Pos, ast.AST, error) {
	return p.Concat(
		func(asts []ast.AST) ast.AST {
			return &ast.Equal{LHS: asts[0], RHS: asts[3]}
		},
		p.Expr,
		p.Skip(kind.Equal),
		p.Skip(kind.Equal),
		p.Cond,
	)(at)
}

func (p *Parser) Expr(at Pos) (Pos, ast.AST, error) {
	return p.Select(p.Add, p.Sub, p.Term)(at)
}

func (p *Parser) Add(at Pos) (Pos, ast.AST, error) {
	return p.Concat(
		func(asts []ast.AST) ast.AST {
			return &ast.Add{LHS: asts[0], RHS: asts[2]}
		},
		p.Term,
		p.Skip(kind.Plus),
		p.Expr,
	)(at)
}

func (p *Parser) Sub(at Pos) (Pos, ast.AST, error) {
	return p.Concat(
		func(asts []ast.AST) ast.AST {
			return &ast.Sub{LHS: asts[0], RHS: asts[2]}
		},
		p.Term,
		p.Skip(kind.Minus),
		p.Expr,
	)(at)
}

func (p *Parser) Term(at Pos) (Pos, ast.AST, error) {
	return p.Select(p.Mul, p.Div, p.Res)(at)
}

func (p *Parser) Mul(at Pos) (Pos, ast.AST, error) {
	return p.Concat(
		func(asts []ast.AST) ast.AST {
			return &ast.Mul{LHS: asts[0], RHS: asts[2]}
		},
		p.Res,
		p.Skip(kind.Multiply),
		p.Term,
	)(at)
}

func (p *Parser) Div(at Pos) (Pos, ast.AST, error) {
	m := func(asts []ast.AST) ast.AST {
		return &ast.Div{LHS: asts[0], RHS: asts[2]}
	}
	return p.Concat(m, p.Res, p.Skip(kind.Divide), p.Term)(at)
}

func (p *Parser) Res(at Pos) (Pos, ast.AST, error) {
	return p.Select(p.If, p.Clause, p.Variable, p.Integer)(at)
}

func (p *Parser) Clause(at Pos) (Pos, ast.AST, error) {
	return p.Concat(
		func(asts []ast.AST) ast.AST { return asts[1] },
		p.Skip(kind.LeftParen),
		p.Cond,
		p.Skip(kind.RightParen),
	)(at)
}

func (p *Parser) If(at Pos) (Pos, ast.AST, error) {
	return p.Concat(
		func(asts []ast.AST) ast.AST {
			return &ast.If{
				Cond: asts[1],
				Then: asts[3],
				Else: asts[7],
			}
		},
		p.Skip(kind.If),
		p.Execute,
		p.Skip(kind.LeftCurly),
		p.Execute,
		p.Skip(kind.RightCurly),
		p.Skip(kind.Else),
		p.Skip(kind.LeftCurly),
		p.Execute,
		p.Skip(kind.RightCurly),
	)(at)
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

func (p *Parser) Variable(at Pos) (Pos, ast.AST, error) {
	nx, t := p.Consume(kind.Identifier, at)
	if t == nil {
		return at, nil, errors.New("invalid token")
	}
	return nx, &ast.Variable{VarName: t.Str}, nil
}

func (p *Parser) String(at Pos) (Pos, ast.AST, error) {
	nx, t := p.Consume(kind.String, at)
	if t == nil {
		return at, nil, errors.New("invalid token")
	}
	word := t.Str[1 : len(t.Str)-1]
	return nx, &ast.String{Word: word}, nil
}

func Parse(tks []*token.Token) (ast.AST, error) {
	parser := Parser{tokens: tks, cache: Cache{}}
	_, ast, err := parser.Root(0)
	return ast, err
}
