package ast

import "github.com/yuniruyuni/lang/ir"

type Call struct {
	Result  Reg
	CondReg Reg

	// for `Name(x, y, z, )`,
	FuncName AST
	Params   AST
}

func (s *Call) Name() string {
	return ""
}

func (s *Call) ResultReg() Reg {
	return s.Result
}

func (s *Call) ResultLabel() Label {
	return 0
}

func (s *Call) GenHeader() ir.IR {
	return ""
}

func (s *Call) GenBody(g *Gen) ir.IR {
	return ""
}

func (s *Call) GenPrinter() ir.IR {
	return GenIntPrinter(s.Result)
}
