package ast

import "github.com/yuniruyuni/lang/ir"

type Params struct {
	Result  Reg
	CondReg Reg

	// for `x, y, z,`
	Vars []AST
}

func (s *Params) Name() string {
	return ""
}

func (s *Params) ResultReg() Reg {
	return s.Result
}

func (s *Params) ResultLabel() Label {
	return 0
}

func (s *Params) GenHeader() ir.IR {
	return ""
}

func (s *Params) GenBody(g *Gen) ir.IR {
	return ""
}

func (s *Params) GenPrinter() ir.IR {
	return GenIntPrinter(s.Result)
}
