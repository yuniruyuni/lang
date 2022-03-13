package ast

import (
	"fmt"
)

type String struct {
	Word string
}

func (s *String) ResultReg() Reg {
	return 0
}

func (s *String) AcquireReg(g *Gen) {
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

func (nd *String) GenHeader() IR {
	template := `@.%s = private unnamed_addr constant [%d x i8] c"%s\00", align 1` + "\n"

	n := nd.Name()
	w := nd.Word
	l := nd.WordLen()

	return IR(fmt.Sprintf(template, n, l, w))
}

func (nd *String) GenBody() IR {
	template := `call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([%d x i8], [%d x i8]* @.%s, i64 0, i64 0))`

	n := nd.Name()
	l := nd.WordLen()

	return IR(fmt.Sprintf(template, l, l, n))
}

func (s *String) GenPrinter() IR {
	return ""
}
