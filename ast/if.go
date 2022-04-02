package ast

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

	jumpBody := IR(`
		%%%d = icmp ne i32 %%%d, 0
		br i1 %%%d, label %%label.%d, label %%label.%d
	`).Expand(
		s.CondReg,
		s.Cond.ResultReg(),
		s.CondReg,
		s.ThenLabel,
		s.ElseLabel,
	)

	phiBody := IR(`
		%%%d = phi i32 [ %%%d, %%label.%d ], [ %%%d, %%label.%d ]
	`).Expand(
		s.ResultReg(),
		s.Then.ResultReg(),
		s.Then.ResultLabel(),
		s.Else.ResultReg(),
		s.Else.ResultLabel(),
	)

	return Concat(
		condBody,
		jumpBody,
		IR("label.%d:").Expand(s.ThenLabel),
		thenBody,
		IR("\t\tbr label %%label.%d").Expand(s.PhiLabel),
		IR("label.%d:").Expand(s.ElseLabel),
		elseBody,
		IR("\t\tbr label %%label.%d\n").Expand(s.PhiLabel),
		IR("label.%d:\n").Expand(s.PhiLabel),
		phiBody,
	)
}

func (s *If) GenPrinter() IR {
	n := "intfmt"
	l := 4
	v := s.Result

	return IR(`call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([%d x i8], [%d x i8]* @.%s, i64 0, i64 0), i32 %%%d)`).
		Expand(l, l, n, v)
}
