package ast

import "github.com/yuniruyuni/lang/ir"

type Less struct {
	TmpReg Reg
	Result Reg
	// for `x < y`,
	LHS AST // x
	RHS AST // y
}

func (s *Less) ResultReg() Reg {
	return s.Result
}

func (s *Less) ResultLabel() Label {
	return s.RHS.ResultLabel()
}

func (s *Less) GenHeader() ir.IR {
	return s.LHS.GenHeader() + s.RHS.GenHeader()
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

func (s *Less) GenPrinter() ir.IR {
	n := "intfmt"
	l := 4
	v := s.Result

	return ir.IR(`call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([%d x i8], [%d x i8]* @.%s, i64 0, i64 0), i32 %%%d)`).
		Expand(l, l, n, v)
}
