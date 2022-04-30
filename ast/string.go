package ast

import (
	"fmt"

	"github.com/yuniruyuni/lang/ir"
)

type String struct {
	NamePostfix Constant
	Word        string
}

func (nd *String) Name() Name {
	return Name(fmt.Sprintf("str.%d", nd.NamePostfix))
}

func (s *String) Type() Type {
	return "i8*"
}

func (nd *String) ResultReg() Reg {
	return 0
}

func (nd *String) ResultLabel() Label {
	return 0
}

const (
	nullCharSize = 1
)

func (nd *String) WordLen() int {
	return len(nd.Word) + nullCharSize
}

func (nd *String) GenHeader(g *Gen) ir.IR {
	nd.NamePostfix = g.NextConstant()

	n := nd.Name()
	w := nd.Word
	l := nd.WordLen()

	return ir.IR(`@.%s = private unnamed_addr constant [%d x i8] c"%s\00", align 1`).
		Expand(n, l, w)
}

func (nd *String) GenBody(g *Gen) ir.IR {
	return ""
}

func (nd *String) GenArg() ir.IR {
	n := nd.Name()
	l := nd.WordLen()
	return ir.IR(`
		i8* getelementptr inbounds (
			[%d x i8],
			[%d x i8]* @.%s,
			i64 0,
			i64 0
		)
	`).Expand(
		l, l, n,
	)
}

func (nd *String) GenPrinter() ir.IR {
	return ir.IR(`call i32 (i8*, ...) @printf(%s)`).Expand(nd.GenArg())
}
