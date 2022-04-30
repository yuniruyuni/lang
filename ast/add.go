package ast

import "github.com/yuniruyuni/lang/ir"

type Add struct {
	Result Reg
	// for `x + y`,
	LHS AST // x
	RHS AST // y
}

func (s *Add) Name() Name {
	return ""
}

func (s *Add) ResultReg() Reg {
	return s.Result
}

func (s *Add) ResultLabel() Label {
	return s.RHS.ResultLabel()
}

func (s *Add) GenHeader(g *Gen) ir.IR {
	return s.LHS.GenHeader(g) + s.RHS.GenHeader(g)
}

func (s *Add) GenBody(g *Gen) ir.IR {
	lhsBody := s.LHS.GenBody(g)
	rhsBody := s.RHS.GenBody(g)
	s.Result = g.NextReg()

	body := ir.IR(`%%%d = add i32 %%%d, %%%d`).
		Expand(s.Result, s.LHS.ResultReg(), s.RHS.ResultReg())

	return ir.Concat(lhsBody, rhsBody, body)
}

func (s *Add) GenArg() ir.IR {
	return ir.IR(`i32 %%%d`).Expand(s.ResultReg())
}

func (s *Add) GenPrinter() ir.IR {
	return GenIntPrinter(s.Result)
}
