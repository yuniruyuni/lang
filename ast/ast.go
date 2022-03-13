package ast

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
