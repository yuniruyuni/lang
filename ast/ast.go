package ast

import "fmt"

type IR string

type AST interface {
	GenHeader() IR
	GenBody() IR
}

type Sum struct {
	// for `x + y`,
	LHS AST // x
	RHS AST // y
}

func (s *Sum) GenHeader() IR {
	return s.LHS.GenHeader() + s.RHS.GenHeader()
}

func (s *Sum) Gen() IR {
	return ``
}

type String struct {
	Word string
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
	template := `%%2 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([%d x i8], [%d x i8]* @.%s, i64 0, i64 0))`

	n := nd.Name()
	l := nd.WordLen()

	return IR(fmt.Sprintf(template, l, l, n))
}

type Integer struct {
	Value int
}

func (nd *Integer) GenHeader() IR {
	return "\n"
}

func (nd *Integer) GenBody() IR {
	template := `%%2 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([%d x i8], [%d x i8]* @.%s, i64 0, i64 0), i32 %d)`

	n := "intfmt"
	l := 4
	v := nd.Value

	return IR(fmt.Sprintf(template, l, l, n, v))
}
