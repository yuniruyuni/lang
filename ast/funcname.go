package ast

import (
	"github.com/yuniruyuni/lang/ir"
)

type FuncName struct {
	FuncName string
}

func (s *FuncName) Name() string {
	return s.FuncName
}

func (s *FuncName) ResultReg() Reg {
	return 0
}

func (s *FuncName) ResultLabel() Label {
	return 0
}

func (s *FuncName) GenHeader() ir.IR {
	return ""
}

func (s *FuncName) GenBody(g *Gen) ir.IR {
	return ""
}

func (s *FuncName) GenPrinter() ir.IR {
	return ""
}
