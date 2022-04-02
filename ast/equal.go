package ast

import "github.com/yuniruyuni/lang/ir"

type Equal struct {
	TmpReg Reg
	Result Reg
	// for `x == y`,
	LHS AST // x
	RHS AST // y
}

func (s *Equal) ResultReg() Reg {
	return s.Result
}

func (s *Equal) ResultLabel() Label {
	return s.RHS.ResultLabel()
}

func (s *Equal) GenHeader() ir.IR {
	return s.LHS.GenHeader() + s.RHS.GenHeader()
}

func (s *Equal) GenBody(g *Gen) ir.IR {
	lhsBody := s.LHS.GenBody(g)
	rhsBody := s.RHS.GenBody(g)
	s.TmpReg = g.NextReg()
	s.Result = g.NextReg()

	body := ir.IR(`
		%%%d = icmp eq i32 %%%d, %%%d
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

func (s *Equal) GenPrinter() ir.IR {
	return GenIntPrinter(s.Result)
}
