package ast

import "github.com/yuniruyuni/lang/ir"

type While struct {
	Result  Reg
	CondReg Reg

	// for `while <Cond> { <Proc> }`,
	Cond AST
	Proc AST

	EntryLabel Label
	TryLabel   Label
	ProcLabel  Label
	EndLabel   Label
}

func (s *While) Name() string {
	return ""
}

func (s *While) ResultReg() Reg {
	return s.Result
}

func (s *While) ResultLabel() Label {
	return s.EndLabel
}

func (s *While) GenHeader() ir.IR {
	return s.Cond.GenHeader() + s.Proc.GenHeader()
}

func (s *While) GenBody(g *Gen) ir.IR {
	s.EntryLabel = g.NextLabel()
	s.TryLabel = g.NextLabel()
	s.Result = g.NextReg()
	condBody := s.Cond.GenBody(g)
	s.CondReg = g.NextReg()
	s.ProcLabel = g.NextLabel()
	procBody := s.Proc.GenBody(g)
	s.EndLabel = g.NextLabel()

	return ir.IR(`
		; ------- entry
		br label %%label.%d

		; ------- condition
		label.%d:
		%%%d = phi i32 [ %%%d, %%label.%d ]

		%s

		%%%d = icmp ne i32 %%%d, 0
		br i1 %%%d, label %%label.%d, label %%label.%d

		; ------- loop clause
		label.%d:
		%s
		br label %%label.%d

		; ------- label for ending loop
		label.%d:
	`).Expand(
		s.TryLabel,
		s.TryLabel,
		s.Result, s.Proc.ResultReg(), s.ProcLabel,
		condBody,
		s.CondReg, s.Cond.ResultReg(),
		s.CondReg, s.ProcLabel, s.EndLabel,
		s.ProcLabel,
		procBody,
		s.TryLabel,
		s.EndLabel,
	)
}

func (s *While) GenPrinter() ir.IR {
	return s.Proc.GenPrinter()
}
