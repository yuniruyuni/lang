package ast

import (
	"github.com/yuniruyuni/lang/ir"
)

type Let struct {
	Result Reg
	// for `let x = y`,
	LHS AST // x
	RHS AST // y
}

func (s *Let) Name() string {
	return s.LHS.Name()
}

func (s *Let) ResultReg() Reg {
	return s.Result
}

func (s *Let) ResultLabel() Label {
	return s.RHS.ResultLabel()
}

func (s *Let) GenHeader() ir.IR {
	return s.RHS.GenHeader()
}

func (s *Let) GenBody(g *Gen) ir.IR {
	rhsBody := s.RHS.GenBody(g)
	s.Result = g.NextReg()

	body := ir.IR(`
		; let X = Y
		%%%s = alloca i32, align 4
		store i32 %%%d, i32* %%%s, align 4
		%%%d = load i32, i32* %%%s, align 4
	`).
		Expand(
			s.Name(),
			s.RHS.ResultReg(), s.Name(),
			s.Result, s.Name(),
		)

	return ir.Concat(rhsBody, body)
}

func (s *Let) GenPrinter() ir.IR {
	return GenIntPrinter(s.Result)
}
