package ast

import "github.com/yuniruyuni/lang/ir"

type Sub struct {
	Result Reg
	// for `x - y`,
	LHS AST // x
	RHS AST // y
}

func (nd *Sub) Name() string {
	return ""
}

func (s *Sub) ResultReg() Reg {
	return s.Result
}

func (s *Sub) ResultLabel() Label {
	return s.RHS.ResultLabel()
}

func (s *Sub) GenHeader() ir.IR {
	return s.LHS.GenHeader() + s.RHS.GenHeader()
}

func (s *Sub) GenBody(g *Gen) ir.IR {
	lhsBody := s.LHS.GenBody(g)
	rhsBody := s.RHS.GenBody(g)
	s.Result = g.NextReg()

	body := ir.IR(`%%%d = sub i32 %%%d, %%%d`).
		Expand(s.Result, s.LHS.ResultReg(), s.RHS.ResultReg())

	return ir.Concat(lhsBody, rhsBody, body)
}

func (s *Sub) GenPrinter() ir.IR {
	return GenIntPrinter(s.Result)
}
