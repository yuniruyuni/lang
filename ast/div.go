package ast

import "github.com/yuniruyuni/lang/ir"

type Div struct {
	Result Reg
	// for `x / y`,
	LHS AST // x
	RHS AST // y
}

func (s *Div) Name() Name {
	return ""
}

func (s *Div) ResultReg() Reg {
	return s.Result
}

func (s *Div) ResultLabel() Label {
	return s.RHS.ResultLabel()
}

func (s *Div) GenHeader(g *Gen) ir.IR {
	return s.LHS.GenHeader(g) + s.RHS.GenHeader(g)
}

func (s *Div) GenBody(g *Gen) ir.IR {
	lhsBody := s.LHS.GenBody(g)
	rhsBody := s.RHS.GenBody(g)
	s.Result = g.NextReg()

	body := ir.IR(`%%%d = sdiv i32 %%%d, %%%d`).
		Expand(s.Result, s.LHS.ResultReg(), s.RHS.ResultReg())

	return ir.Concat(lhsBody, rhsBody, body)
}

func (s *Div) GenArg() ir.IR {
	return ir.IR(`i32 %%%d`).Expand(s.ResultReg())
}

func (s *Div) GenPrinter() ir.IR {
	return GenIntPrinter(s.Result)
}
