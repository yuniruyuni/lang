package parse

import (
	"errors"
	"reflect"

	"github.com/yuniruyuni/lang/ast"
)

// CachedCall calls f() and cache the result if not error.
// NOTE: This function just can handle named-functions.
func (p *Parser) CachedCall(f NonTerminal, at Pos) (Pos, ast.AST, error) {
	ptr := reflect.ValueOf(f).Pointer()

	key := Key{Ptr: ptr, At: at}
	res, ok := p.cache[key]
	if ok {
		return res.Pos, res.Ast, nil
	}

	nx, parsed, err := f(at)
	if err == nil {
		p.cache[key] = &Result{Pos: nx, Ast: parsed}
	}
	return nx, parsed, err
}

// Select combines some target NonTerminals into single NonTerminal.
// This new NonTerminal checks if targets match current tokens and
// returns first maching NonTerminal result.
// If there is no maching NonTerminal, It returns an "invalid tokens" error.
func (p *Parser) Select(cands ...NonTerminal) NonTerminal {
	return func(at Pos) (Pos, ast.AST, error) {
		for _, cand := range cands {
			nx, parsed, err := p.CachedCall(cand, at)
			if err == nil {
				return nx, parsed, nil
			}
		}
		return at, nil, errors.New("invalid tokens")
	}
}

// Merger merges some of ASTs into single AST node.
type Merger func([]ast.AST) ast.AST

// Concat combines NonTerminals sequence into single NonTerminal.
// This new NonTerminal checks if sequence match from current tokens and
// call Merger by matched ASTs then the Merger's result return.
// If there is a non-matched NonTerminal,
// It returns the error of the NonTerminal.
func (p *Parser) Concat(m Merger, cands ...NonTerminal) NonTerminal {
	return func(at Pos) (Pos, ast.AST, error) {
		asts := make([]ast.AST, 0, len(cands))

		nx := at
		for _, cand := range cands {
			var parsed ast.AST
			var err error
			nx, parsed, err = p.CachedCall(cand, nx)
			if err != nil {
				return at, nil, err
			}
			asts = append(asts, parsed)
		}
		return nx, m(asts), nil
	}
}

// Many takes a (empty/non-empty) sequence of NonTerminal and it makes single NonTerminal.
func (p *Parser) Many(m Merger, cand NonTerminal) NonTerminal {
	return func(at Pos) (Pos, ast.AST, error) {
		asts := make([]ast.AST, 0)

		nx := at
		for {
			var err error
			var parsed ast.AST
			nx, parsed, err = p.CachedCall(cand, nx)
			if err != nil {
				return nx, m(asts), nil
			}
			asts = append(asts, parsed)
		}
	}
}
