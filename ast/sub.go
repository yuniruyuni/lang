package ast

import (
	"fmt"
	"strings"
)

type Sub struct {
	Result Reg
	// for `x - y`,
	LHS AST // x
	RHS AST // y
}

func (s *Sub) ResultReg() Reg {
	return s.Result
}

func (s *Sub) ResultLabel() Label {
	return s.RHS.ResultLabel()
}

func (s *Sub) GenHeader() IR {
	return s.LHS.GenHeader() + s.RHS.GenHeader()
}

func (s *Sub) GenBody(g *Gen) IR {
	lhsBody := s.LHS.GenBody(g)
	rhsBody := s.RHS.GenBody(g)
	s.Result = g.NextReg()

	tmpl := `
		%%%d = sub i32 %%%d, %%%d
	`
	body := fmt.Sprintf(
		tmpl,
		s.Result,
		s.LHS.ResultReg(),
		s.RHS.ResultReg(),
	)

	bodies := []string{
		string(lhsBody),
		string(rhsBody),
		string(body),
	}

	return IR(strings.Join(bodies, "\n"))
}

func (s *Sub) GenPrinter() IR {
	tmpl := `
		call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([%d x i8], [%d x i8]* @.%s, i64 0, i64 0), i32 %%%d)
	`

	n := "intfmt"
	l := 4
	v := s.Result

	return IR(fmt.Sprintf(tmpl, l, l, n, v))
}
