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
	return ""
}

func (s *Definitions) GenBody(g *Gen) ir.IR {
	return ""
}

func (s *Definitions) GenArg() ir.IR {
	return ir.IR(`i32 %%%d`).Expand(s.ResultReg())
}

func (s *Definitions) GenPrinter() ir.IR {
	return GenIntPrinter(s.Result)
}
