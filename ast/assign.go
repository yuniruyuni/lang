package ast

import (
	"github.com/yuniruyuni/lang/ir"
)

type Assign struct {
	Result Reg
	// for `x = y`,
	LHS AST // x
	RHS AST // y
}

func (s *Assign) Name() Name {
	return s.LHS.Name()
}

func (s *Assign) Type() Type {
	return "i32"
}

func (s *Assign) ResultReg() Reg {
	return s.Result
}

func (s *Assign) ResultLabel() Label {
	return s.RHS.ResultLabel()
}

func (s *Assign) GenHeader(g *Gen) ir.IR {
	return s.RHS.GenHeader(g)
}

func (s *Assign) GenBody(g *Gen) ir.IR {
	rhsBody := s.RHS.GenBody(g)
	s.Result = g.NextReg()

	body := ir.IR(`
		; X = Y
		store i32 %%%d, i32* %%%s, align 4
		%%%d = load i32, i32* %%%s, align 4
	`).
		Expand(
			s.RHS.ResultReg(), s.Name(),
			s.Result, s.Name(),
		)

	return ir.Concat(rhsBody, body)
}

func (s *Assign) GenArg() ir.IR {
	return ir.IR(`i32 %%%d`).Expand(s.ResultReg())
}

func (s *Assign) GenPrinter() ir.IR {
	return GenIntPrinter(s.Result)
}
