package ast

import (
	"github.com/yuniruyuni/lang/ir"
)

type FuncName struct {
	FuncName Name
}

func (s *FuncName) Name() Name {
	return s.FuncName
}

func (s *FuncName) Type() Type {
	return ""
}

func (s *FuncName) ResultReg() Reg {
	return 0
}

func (s *FuncName) ResultLabel() Label {
	return 0
}

func (s *FuncName) GenHeader(g *Gen) ir.IR {
	return ""
}

func (s *FuncName) GenBody(g *Gen) ir.IR {
	return ""
}

func (s *FuncName) GenArg() ir.IR {
	return ir.IR(`i32 %%%d`).Expand(s.ResultReg())
}

func (s *FuncName) GenPrinter() ir.IR {
	return ""
}
