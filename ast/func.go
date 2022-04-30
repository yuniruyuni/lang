package ast

import (
	"github.com/yuniruyuni/lang/ir"
)

type Func struct {
	FuncName AST
	Params   AST
	Execute  AST
}

func (s *Func) Name() Name {
	return s.FuncName.Name()
}

func (s *Func) ResultReg() Reg {
	return 0
}

func (s *Func) ResultLabel() Label {
	return 0
}

func (s *Func) GenHeader(g *Gen) ir.IR {
	return ""
}

func (s *Func) GenBody(g *Gen) ir.IR {
	return ""
}

func (s *Func) GenArg() ir.IR {
	return ir.IR(`i32 %%%d`).Expand(s.ResultReg())
}

func (s *Func) GenPrinter() ir.IR {
	return ""
}
