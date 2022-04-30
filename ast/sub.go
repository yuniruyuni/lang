package ast

import "github.com/yuniruyuni/lang/ir"

type Sub struct {
	Result Reg
	// for `x - y`,
	LHS AST // x
	RHS AST // y
}

func (nd *Sub) Name() Name {
	return ""
}

func (s *Sub) ResultReg() Reg {
	return s.Result
}

func (s *Sub) ResultLabel() Label {
	return s.RHS.ResultLabel()
}

func (s *Sub) GenHeader(g *Gen) ir.IR {
	return s.LHS.GenHeader(g) + s.RHS.GenHeader(g)
}

func (s *Sub) GenBody(g *Gen) ir.IR {
	lhsBody := s.LHS.GenBody(g)
	rhsBody := s.RHS.GenBody(g)
	s.Result = g.NextReg()

	body := ir.IR(`%%%d = sub i32 %%%d, %%%d`).
		Expand(s.Result, s.LHS.ResultReg(), s.RHS.ResultReg())

	return ir.Concat(lhsBody, rhsBody, body)
}

func (s *Sub) GenArg() ir.IR {
	return ir.IR(`i32 %%%d`).Expand(s.ResultReg())
}

func (s *Sub) GenPrinter() ir.IR {
	return GenIntPrinter(s.Result)
}
