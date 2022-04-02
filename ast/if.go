package ast

import (
	"fmt"
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

func (s *If) GenHeader() IR {
	return s.Cond.GenHeader() + s.Then.GenHeader() + s.Else.GenHeader()
}

func (s *If) GenBody(g *Gen) IR {
	condBody := s.Cond.GenBody(g)
	s.CondReg = g.NextReg()
	s.ThenLabel = g.NextLabel()
	thenBody := s.Then.GenBody(g)
	s.ElseLabel = g.NextLabel()
	elseBody := s.Else.GenBody(g)
	s.PhiLabel = g.NextLabel()
	s.Result = g.NextReg()

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

	return Concat(
		condBody,
		IR(jumpBody),
		IR(fmt.Sprintf("label.%d:\n", s.ThenLabel)),
		thenBody,
		IR(fmt.Sprintf("\t\tbr label %%label.%d\n", s.PhiLabel)),
		IR(fmt.Sprintf("label.%d:\n", s.ElseLabel)),
		elseBody,
		IR(fmt.Sprintf("\t\tbr label %%label.%d\n", s.PhiLabel)),
		IR(fmt.Sprintf("label.%d:\n", s.PhiLabel)),
		IR(phiBody),
	)
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
