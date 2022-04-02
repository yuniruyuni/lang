package ast

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

type AST interface {
	ResultReg() Reg
	ResultLabel() Label

	GenHeader() IR
	GenBody(g *Gen) IR
	GenPrinter() IR
}
