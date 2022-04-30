package ast

import (
	"github.com/yuniruyuni/lang/ir"
)

type Func struct {
	FuncName AST
	Params   AST
	Execute  AST
}

func (s *Func) Name() Name {
	return s.FuncName.Name()
}

func (s *Func) ResultReg() Reg {
	return 0
}

func (s *Func) ResultLabel() Label {
	return 0
}

func (s *Func) GenHeader(g *Gen) ir.IR {
	return s.Execute.GenHeader(g)
}

func (s *Func) GenBody(g *Gen) ir.IR {
	name := s.Name()
	params := s.Params.GenBody(g)
	body := s.Execute.GenBody(g)

	return ir.IR(`
		define i32 @%s(%s) {
			%s
			ret i32 %%%d
		}
	`).Expand(
		name, params,
		body,
		s.Execute.ResultReg(),
	)
}

func (s *Func) GenArg() ir.IR {
	return ""
}

func (s *Func) GenPrinter() ir.IR {
	return ""
}
