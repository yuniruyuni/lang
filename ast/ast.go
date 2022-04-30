package ast

import (
	"fmt"

	"github.com/yuniruyuni/lang/ir"
)

type Reg int
type Label int
type Constant int

type Name string
type Type []string

type Gen struct {
	reg      Reg
	label    Label
	constant Constant
	types    map[Name]Type
}

func NewGen() *Gen {
	return &Gen{
		reg:      0,
		label:    0,
		constant: 0,
		types:    map[Name]Type{},
	}
}

func (g *Gen) CurConstant() Constant {
	return g.constant
}

func (g *Gen) NextConstant() Constant {
	g.constant += 1
	return g.constant
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

func (g *Gen) RegisterFunc(n Name, t Type) {
	g.types[n] = t
}

func (g *Gen) GetFunc(n Name) (Type, error) {
	t, ok := g.types[n]
	if !ok {
		return nil, fmt.Errorf("Function %s doesn't exist.", n)
	}
	return t, nil
}

type AST interface {
	ResultReg() Reg
	ResultLabel() Label

	Name() Name

	GenHeader(g *Gen) ir.IR
	GenBody(g *Gen) ir.IR
	GenArg() ir.IR
	GenPrinter() ir.IR
}
