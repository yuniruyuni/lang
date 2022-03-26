package ast

import (
	"fmt"
)

type Integer struct {
	Result Reg
	Label  Label
	Alloc  Reg
	Value  int
}

func (nd *Integer) ResultReg() Reg {
	return nd.Result
}

func (s *Integer) ResultLabel() Label {
	return s.Label
}

func (nd *Integer) AcquireReg(g *Gen) {
	nd.Alloc = g.NextReg()
	nd.Result = g.NextReg()
	nd.Label = g.CurLabel()
}

func (nd *Integer) GenHeader() IR {
	return "\n"
}

func (nd *Integer) GenBody() IR {
	template := `
		%%%d = alloca i32, align 4
		store i32 %d, i32* %%%d
		%%%d = load i32, i32* %%%d, align 4
	`

	return IR(fmt.Sprintf(template, nd.Alloc, nd.Value, nd.Alloc, nd.Result, nd.Alloc))
}

func (nd *Integer) GenPrinter() IR {
	tmpl := `
		call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([%d x i8], [%d x i8]* @.%s, i64 0, i64 0), i32 %%%d)
	`

	n := "intfmt"
	l := 4
	v := nd.Result

	return IR(fmt.Sprintf(tmpl, l, l, n, v))
}
