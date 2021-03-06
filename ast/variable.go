package ast

import (
	"github.com/yuniruyuni/lang/ir"
)

type Variable struct {
	Result  Reg
	Label   Label
	VarName Name
}

func (s *Variable) Name() Name {
	return s.VarName
}

func (s *Variable) Type() Type {
	return "i32"
}

func (s *Variable) ResultReg() Reg {
	return s.Result
}

func (s *Variable) ResultLabel() Label {
	return s.Label
}

func (s *Variable) GenHeader(g *Gen) ir.IR {
	return ""
}

func (s *Variable) GenBody(g *Gen) ir.IR {
	s.Result = g.NextReg()
	s.Label = g.CurLabel()

	vartype, err := g.GetVariable(s.Name())
	if err != nil {
		panic(err)
	}

	switch vartype {
	case "i32":
		return ir.IR(`%%%d = add i32 %%%s, 0`).
			Expand(s.Result, s.Name())
	case "i32*":
		return ir.IR(`%%%d = load i32, i32* %%%s, align 4`).
			Expand(s.Result, s.Name())
	default:
		panic("unsupported type has held")
	}
}

func (s *Variable) GenArg() ir.IR {
	return ir.IR(`i32 %%%d`).Expand(s.ResultReg())
}

func (s *Variable) GenPrinter() ir.IR {
	return GenIntPrinter(s.Result)
}
