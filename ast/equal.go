package ast

import (
	"fmt"
	"strings"
)

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

func (s *Equal) GenHeader() IR {
	return s.LHS.GenHeader() + s.RHS.GenHeader()
}

func (s *Equal) GenBody(g *Gen) IR {
	lhsBody := s.LHS.GenBody(g)
	rhsBody := s.RHS.GenBody(g)
	s.TmpReg = g.NextReg()
	s.Result = g.NextReg()

	tmpl := `
		%%%d = icmp eq i32 %%%d, %%%d
		%%%d = zext i1 %%%d to i32
	`
	body := fmt.Sprintf(
		tmpl,
		s.TmpReg,
		s.LHS.ResultReg(),
		s.RHS.ResultReg(),
		s.Result,
		s.TmpReg,
	)

	bodies := []string{
		string(lhsBody),
		string(rhsBody),
		string(body),
	}

	return IR(strings.Join(bodies, "\n"))
}

func (s *Equal) GenPrinter() IR {
	tmpl := `
		call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([%d x i8], [%d x i8]* @.%s, i64 0, i64 0), i32 %%%d)
	`

	n := "intfmt"
	l := 4
	v := s.Result

	return IR(fmt.Sprintf(tmpl, l, l, n, v))
}
