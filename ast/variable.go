package ast

import (
	"github.com/yuniruyuni/lang/ir"
)

type Variable struct {
	VarName string
}

func (s *Variable) Name() string {
	return s.VarName
}

func (s *Variable) ResultReg() Reg {
	// TODO: implement
	return 0
}

func (s *Variable) ResultLabel() Label {
	// TODO: implement
	return 0
}

func (s *Variable) GenHeader() ir.IR {
	// TODO: implement
	return ""
}

func (s *Variable) GenBody(g *Gen) ir.IR {
	// TODO: implement
	return ""
}

func (s *Variable) GenPrinter() ir.IR {
	// TODO: implement
	return ""
}
