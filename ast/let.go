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

func (s *Let) Name() Name {
	return s.LHS.Name()
}

func (s *Let) Type() Type {
	return "i32"
}

func (s *Let) ResultReg() Reg {
	return s.Result
}

func (s *Let) ResultLabel() Label {
	return s.RHS.ResultLabel()
}

func (s *Let) GenHeader(g *Gen) ir.IR {
	return s.RHS.GenHeader(g)
}

func (s *Let) GenBody(g *Gen) ir.IR {
	g.RegisterVariable(s.Name(), "i32*")

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

func (s *Let) GenArg() ir.IR {
	return ir.IR(`i32 %%%d`).Expand(s.ResultReg())
}

func (s *Let) GenPrinter() ir.IR {
	return GenIntPrinter(s.Result)
}
