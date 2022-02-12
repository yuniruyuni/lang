package ast

import (
	"fmt"
	"strings"
)

type Reg int

type Gen struct {
	reg Reg
}

func (g *Gen) CurReg() Reg {
	return g.reg
}

func (g *Gen) NextReg() Reg {
	g.reg += 1
	return g.reg
}

type IR string

type AST interface {
	ResultReg() Reg
	AcquireReg(g *Gen)

	GenHeader() IR
	GenBody() IR
	GenPrinter() IR
}

type Sum struct {
	Result Reg
	// for `x + y`,
	LHS AST // x
	RHS AST // y
}

func (s *Sum) ResultReg() Reg {
	return s.Result
}

func (s *Sum) AcquireReg(g *Gen) {
	s.LHS.AcquireReg(g)
	s.RHS.AcquireReg(g)
	s.Result = g.NextReg()
}

func (s *Sum) GenHeader() IR {
	return s.LHS.GenHeader() + s.RHS.GenHeader()
}

func (s *Sum) GenBody() IR {
	lhsBody := s.LHS.GenBody()
	rhsBody := s.RHS.GenBody()

	tmpl := `
		%%%d = add i32 %%%d, %%%d
	`
	body := fmt.Sprintf(
		tmpl,
		s.Result,
		s.LHS.ResultReg(),
		s.RHS.ResultReg(),
	)

	bodies := []string{
		string(lhsBody),
		string(rhsBody),
		string(body),
	}

	return IR(strings.Join(bodies, "\n"))
}

func (s *Sum) GenPrinter() IR {
	tmpl := `
		call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([%d x i8], [%d x i8]* @.%s, i64 0, i64 0), i32 %%%d)
	`

	n := "intfmt"
	l := 4
	v := s.Result

	return IR(fmt.Sprintf(tmpl, l, l, n, v))
}

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

type Integer struct {
	Result Reg
	Alloc  Reg
	Value  int
}

func (nd *Integer) ResultReg() Reg {
	return nd.Result
}

func (nd *Integer) AcquireReg(g *Gen) {
	nd.Alloc = g.NextReg()
	nd.Result = g.NextReg()
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
