package ast

import "github.com/yuniruyuni/lang/ir"

type Less struct {
	TmpReg Reg
	Result Reg
	// for `x < y`,
	LHS AST // x
	RHS AST // y
}

func (s *Less) Name() Name {
	return ""
}

func (s *Less) ResultReg() Reg {
	return s.Result
}

func (s *Less) ResultLabel() Label {
	return s.RHS.ResultLabel()
}

func (s *Less) GenHeader(g *Gen) ir.IR {
	return s.LHS.GenHeader(g) + s.RHS.GenHeader(g)
}

func (s *Less) GenBody(g *Gen) ir.IR {
	lhsBody := s.LHS.GenBody(g)
	rhsBody := s.RHS.GenBody(g)
	s.TmpReg = g.NextReg()
	s.Result = g.NextReg()

	body := ir.IR(`
		%%%d = icmp slt i32 %%%d, %%%d
		%%%d = zext i1 %%%d to i32
	`).Expand(
		s.TmpReg,
		s.LHS.ResultReg(),
		s.RHS.ResultReg(),
		s.Result,
		s.TmpReg,
	)

	return ir.Concat(lhsBody, rhsBody, body)
}

func (s *Less) GenArg() ir.IR {
	return ir.IR(`i32 %%%d`).Expand(s.ResultReg())
}

func (s *Less) GenPrinter() ir.IR {
	return GenIntPrinter(s.Result)
}
