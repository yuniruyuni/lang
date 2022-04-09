package ast

import "github.com/yuniruyuni/lang/ir"

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

func (s *If) Name() string {
	return ""
}

func (s *If) ResultReg() Reg {
	return s.Result
}

func (s *If) ResultLabel() Label {
	return s.PhiLabel
}

func (s *If) GenHeader() ir.IR {
	return s.Cond.GenHeader() + s.Then.GenHeader() + s.Else.GenHeader()
}

func (s *If) GenBody(g *Gen) ir.IR {
	condBody := s.Cond.GenBody(g)
	s.CondReg = g.NextReg()
	s.ThenLabel = g.NextLabel()
	thenBody := s.Then.GenBody(g)
	s.ElseLabel = g.NextLabel()
	elseBody := s.Else.GenBody(g)
	s.PhiLabel = g.NextLabel()
	s.Result = g.NextReg()

	return ir.IR(`
		; ------- start if condition
		%s

		; ------- check the condition meets or not
		%%%d = icmp ne i32 %%%d, 0
		br i1 %%%d, label %%label.%d, label %%label.%d

		; ------- then clause
		label.%d:
		%s
		br label %%label.%d

		; ------- else clause
		label.%d:
		%s
		br label %%label.%d

		; ------- phi label for an if expression
		label.%d:
		%%%d = phi i32 [ %%%d, %%label.%d ], [ %%%d, %%label.%d ]
	`).Expand(
		condBody,
		s.CondReg, s.Cond.ResultReg(),
		s.CondReg, s.ThenLabel, s.ElseLabel,
		s.ThenLabel,
		thenBody,
		s.PhiLabel, s.ElseLabel,
		elseBody,
		s.PhiLabel, s.PhiLabel,
		s.ResultReg(),
		s.Then.ResultReg(),
		s.Then.ResultLabel(),
		s.Else.ResultReg(),
		s.Else.ResultLabel(),
	)
}

func (s *If) GenPrinter() ir.IR {
	return GenIntPrinter(s.Result)
}
