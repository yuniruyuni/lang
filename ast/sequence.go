package ast

import "github.com/yuniruyuni/lang/ir"

type Sequence struct {
	Result Reg
	// for `x; y`,
	LHS AST // x
	RHS AST // y
}

func (s *Sequence) Name() string {
	return ""
}

func (s *Sequence) ResultReg() Reg {
	return s.Result
}

func (s *Sequence) ResultLabel() Label {
	return s.RHS.ResultLabel()
}

func (s *Sequence) GenHeader() ir.IR {
	return s.LHS.GenHeader() + s.RHS.GenHeader()
}

func (s *Sequence) GenBody(g *Gen) ir.IR {
	lhsBody := s.LHS.GenBody(g)
	rhsBody := s.RHS.GenBody(g)
	s.Result = s.RHS.ResultReg()

	return ir.Concat(lhsBody, rhsBody)
}

func (s *Sequence) GenPrinter() ir.IR {
	return s.RHS.GenPrinter()
}
