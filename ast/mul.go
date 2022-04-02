package ast

import "github.com/yuniruyuni/lang/ir"

type Mul struct {
	Result Reg
	// for `x * y`,
	LHS AST // x
	RHS AST // y
}

func (s *Mul) ResultReg() Reg {
	return s.Result
}

func (s *Mul) ResultLabel() Label {
	return s.RHS.ResultLabel()
}

func (s *Mul) GenHeader() ir.IR {
	return s.LHS.GenHeader() + s.RHS.GenHeader()
}

func (s *Mul) GenBody(g *Gen) ir.IR {
	lhsBody := s.LHS.GenBody(g)
	rhsBody := s.RHS.GenBody(g)
	s.Result = g.NextReg()

	body := ir.IR(`%%%d = mul i32 %%%d, %%%d`).
		Expand(s.Result, s.LHS.ResultReg(), s.RHS.ResultReg())

	return ir.Concat(lhsBody, rhsBody, body)
}

func (s *Mul) GenPrinter() ir.IR {
	n := "intfmt"
	l := 4
	v := s.Result

	return ir.IR(`call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([%d x i8], [%d x i8]* @.%s, i64 0, i64 0), i32 %%%d)`).
		Expand(l, l, n, v)
}
