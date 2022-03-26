package ast

import (
	"fmt"
	"strings"
)

type If struct {
	Result  Reg
	CondReg Reg
	// for `if <Cond> { <Then> } else { <Else> }`,
	Cond AST
	Then AST
	Else AST

	ThenLabel Label
	ElseLabel Label
	PhiLabel  Label
}

func (s *If) ResultReg() Reg {
	return s.Result
}

func (s *If) ResultLabel() Label {
	return s.PhiLabel
}

func (s *If) AcquireReg(g *Gen) {
	s.Cond.AcquireReg(g)
	s.CondReg = g.NextReg()
	s.ThenLabel = g.NextLabel()
	s.Then.AcquireReg(g)
	s.ElseLabel = g.NextLabel()
	s.Else.AcquireReg(g)
	s.PhiLabel = g.NextLabel()
	s.Result = g.NextReg()
}

func (s *If) GenHeader() IR {
	return s.Cond.GenHeader() + s.Then.GenHeader() + s.Else.GenHeader()
}

func (s *If) GenBody() IR {
	condBody := s.Cond.GenBody()
	thenBody := s.Then.GenBody()
	elseBody := s.Else.GenBody()

	jumpTmpl := `
		%%%d = icmp ne i32 %%%d, 0
		br i1 %%%d, label %%label.%d, label %%label.%d
	`
	jumpBody := fmt.Sprintf(
		jumpTmpl,
		s.CondReg,
		s.Cond.ResultReg(),
		s.CondReg,
		s.ThenLabel,
		s.ElseLabel,
	)

	phiTmpl := `
		%%%d = phi i32 [ %%%d, %%label.%d ], [ %%%d, %%label.%d ]
	`
	phiBody := fmt.Sprintf(
		phiTmpl,
		s.ResultReg(),
		s.Then.ResultReg(),
		s.Then.ResultLabel(),
		s.Else.ResultReg(),
		s.Else.ResultLabel(),
	)

	bodies := []string{
		string(condBody),
		string(jumpBody),
		fmt.Sprintf("label.%d:\n", s.ThenLabel),
		string(thenBody),
		fmt.Sprintf("\t\tbr label %%label.%d\n", s.PhiLabel),
		fmt.Sprintf("label.%d:\n", s.ElseLabel),
		string(elseBody),
		fmt.Sprintf("\t\tbr label %%label.%d\n", s.PhiLabel),
		fmt.Sprintf("label.%d:\n", s.PhiLabel),
		string(phiBody),
	}

	return IR(strings.Join(bodies, "\n"))
}

func (s *If) GenPrinter() IR {
	tmpl := `
		call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([%d x i8], [%d x i8]* @.%s, i64 0, i64 0), i32 %%%d)
	`

	n := "intfmt"
	l := 4
	v := s.Result

	return IR(fmt.Sprintf(tmpl, l, l, n, v))
}
