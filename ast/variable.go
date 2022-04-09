package ast

import (
	"github.com/yuniruyuni/lang/ir"
)

type Variable struct {
	Result  Reg
	Label   Label
	VarName string
}

func (s *Variable) Name() string {
	return s.VarName
}

func (s *Variable) ResultReg() Reg {
	return s.Result
}

func (s *Variable) ResultLabel() Label {
	return s.Label
}

func (s *Variable) GenHeader() ir.IR {
	return ""
}

func (s *Variable) GenBody(g *Gen) ir.IR {
	s.Result = g.NextReg()
	s.Label = g.CurLabel()

	return ir.IR(`%%%d = load i32, i32* %%%s, align 4`).
		Expand(s.Result, s.Name())
}

func (s *Variable) GenPrinter() ir.IR {
	return GenIntPrinter(s.Result)
}
