package ast

import "github.com/yuniruyuni/lang/ir"

type Definitions struct {
	Result Reg
	// for all definitions
	Defs []AST
}

func (s *Definitions) Name() Name {
	return ""
}

func (s *Definitions) ResultReg() Reg {
	return s.Result
}

func (s *Definitions) ResultLabel() Label {
	return 0
}

func (s *Definitions) GenHeader(g *Gen) ir.IR {
	headers := make([]ir.IR, 0, len(s.Defs))

	for _, d := range s.Defs {
		headers = append(headers, d.GenHeader(g))
	}

	return ir.Concat(headers...)
}

func (s *Definitions) GenBody(g *Gen) ir.IR {
	bodies := make([]ir.IR, 0, len(s.Defs))

	for _, d := range s.Defs {
		g.ResetReg()
		g.ResetLabel()
		bodies = append(bodies, d.GenBody(g))
	}

	return ir.Concat(bodies...)
}

func (s *Definitions) GenArg() ir.IR {
	return ""
}

func (s *Definitions) GenPrinter() ir.IR {
	return ""
}
