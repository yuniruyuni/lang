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

func (s *String) WordLen() int {
	return len(s.Word) + nullCharSize
}

func (s *String) Name() string {
	return "str" // TODO: determine specific name for each ast node.
}

func (s *String) GenHeader() IR {
	template := `@.%s = private unnamed_addr constant [%d x i8] c"%s\00", align 1` + "\n"

	n := s.Name()
	w := s.Word
	l := s.WordLen()

	return IR(fmt.Sprintf(template, n, l, w))
}

func (s *String) GenBody() IR {
	template := `%%2 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([%d x i8], [%d x i8]* @.%s, i64 0, i64 0))`

	n := s.Name()
	l := s.WordLen()

	return IR(fmt.Sprintf(template, l, l, n))
}

type Integer struct {
}

func (s *Integer) GenHeader() IR {
	return "\n"
}

func (s *Integer) GenBody() IR {
	return ``
}
