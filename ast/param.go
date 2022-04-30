package ast

import (
	"github.com/yuniruyuni/lang/ir"
)

type Param struct {
	Result  Reg
	Label   Label
	VarName Name
}

func (s *Param) Name() Name {
	return s.VarName
}

func (s *Param) Type() Type {
	return "i32"
}

func (s *Param) ResultReg() Reg {
	return s.Result
}

func (s *Param) ResultLabel() Label {
	return s.Label
}

func (s *Param) GenHeader(g *Gen) ir.IR {
	return ""
}

func (s *Param) GenBody(g *Gen) ir.IR {
	g.RegisterVariable(s.Name(), "i32")
	return ir.IR(`i32 %%%s`).Expand(s.Name())
}

func (s *Param) GenArg() ir.IR {
	return ""
}

func (s *Param) GenPrinter() ir.IR {
	return GenIntPrinter(s.Result)
}
