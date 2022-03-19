package ast

import (
	"fmt"
	"strings"
)

type Div struct {
	Result Reg
	// for `x / y`,
	LHS AST // x
	RHS AST // y
}

func (s *Div) ResultReg() Reg {
	return s.Result
}

func (s *Div) AcquireReg(g *Gen) {
	s.LHS.AcquireReg(g)
	s.RHS.AcquireReg(g)
	s.Result = g.NextReg()
}

func (s *Div) GenHeader() IR {
	return s.LHS.GenHeader() + s.RHS.GenHeader()
}

func (s *Div) GenBody() IR {
	lhsBody := s.LHS.GenBody()
	rhsBody := s.RHS.GenBody()

	tmpl := `
		%%%d = sdiv i32 %%%d, %%%d
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

func (s *Div) GenPrinter() IR {
	tmpl := `
		call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([%d x i8], [%d x i8]* @.%s, i64 0, i64 0), i32 %%%d)
	`

	n := "intfmt"
	l := 4
	v := s.Result

	return IR(fmt.Sprintf(tmpl, l, l, n, v))
}