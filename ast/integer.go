package ast

import "github.com/yuniruyuni/lang/ir"

type Integer struct {
	Result Reg
	Label  Label
	Alloc  Reg
	Value  int
}

func (s *Integer) Name() Name {
	return ""
}

func (nd *Integer) ResultReg() Reg {
	return nd.Result
}

func (s *Integer) ResultLabel() Label {
	return s.Label
}

func (nd *Integer) GenHeader(g *Gen) ir.IR {
	return "\n"
}

func (nd *Integer) GenBody(g *Gen) ir.IR {
	nd.Alloc = g.NextReg()
	nd.Result = g.NextReg()
	nd.Label = g.CurLabel()

	return ir.IR(`
		%%%d = alloca i32, align 4
		store i32 %d, i32* %%%d
		%%%d = load i32, i32* %%%d, align 4
	`).Expand(
		nd.Alloc, nd.Value, nd.Alloc, nd.Result, nd.Alloc,
	)
}

func (s *Integer) GenArg() ir.IR {
	return ir.IR(`i32 %%%d`).Expand(s.ResultReg())
}

func (nd *Integer) GenPrinter() ir.IR {
	return GenIntPrinter(nd.Result)
}
