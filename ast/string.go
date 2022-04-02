package ast

import "github.com/yuniruyuni/lang/ir"

type String struct {
	Word string
}

func (s *String) ResultReg() Reg {
	return 0
}

func (s *String) ResultLabel() Label {
	return 0
}

const (
	nullCharSize = 1
)

func (nd *String) WordLen() int {
	return len(nd.Word) + nullCharSize
}

func (nd *String) Name() string {
	return "str" // TODO: determine specific name for each ast node.
}

func (nd *String) GenHeader() ir.IR {
	n := nd.Name()
	w := nd.Word
	l := nd.WordLen()

	return ir.IR(`@.%s = private unnamed_addr constant [%d x i8] c"%s\00", align 1`).
		Expand(n, l, w)
}

func (nd *String) GenBody(g *Gen) ir.IR {
	n := nd.Name()
	l := nd.WordLen()

	return ir.IR(`call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([%d x i8], [%d x i8]* @.%s, i64 0, i64 0))`).
		Expand(l, l, n)
}

func (s *String) GenPrinter() ir.IR {
	return ""
}
