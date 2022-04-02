package ast

import (
	"fmt"
	"strings"
)

type Reg int
type Label int

type Gen struct {
	reg   Reg
	label Label
}

func (g *Gen) CurReg() Reg {
	return g.reg
}

func (g *Gen) NextReg() Reg {
	g.reg += 1
	return g.reg
}

func (g *Gen) CurLabel() Label {
	return g.label
}

func (g *Gen) NextLabel() Label {
	g.label += 1
	return g.label
}

type IR string

func (ir IR) Expand(vals ...interface{}) IR {
	return IR(fmt.Sprintf(string(ir), vals...))
}

type IRs []IR

// Concat concatinate all given IRs into single IR.
func Concat(irs ...IR) IR {
	strs := make([]string, 0, len(irs))
	for _, ir := range irs {
		strs = append(strs, string(ir))
	}
	return IR(strings.Join(strs, "\n"))
}

type AST interface {
	ResultReg() Reg
	ResultLabel() Label

	GenHeader() IR
	GenBody(g *Gen) IR
	GenPrinter() IR
}
