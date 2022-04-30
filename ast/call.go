package ast

import (
	"strings"

	"github.com/yuniruyuni/lang/ir"
)

type Call struct {
	Result Reg

	// for `Name(x, y, z, )`,
	FuncName AST
	Args     AST
}

func (s *Call) Name() Name {
	return ""
}

func (s *Call) ResultReg() Reg {
	return s.Result
}

func (s *Call) ResultLabel() Label {
	return 0
}

func (s *Call) GenHeader(g *Gen) ir.IR {
	return s.Args.GenHeader(g)
}

func (s *Call) GenBody(g *Gen) ir.IR {
	argsBody := s.Args.GenBody(g)
	s.Result = g.NextReg()

	types, err := g.GetFunc(Name(s.FuncName.Name()))
	if err != nil {
		panic(err)
	}

	typeDelc := strings.Join(types, ",")

	args := s.Args.GenArg()

	return ir.IR(`
		%s
		%%%d = call i32 (%s) @%s(%s)
	`).Expand(
		argsBody,
		s.Result, typeDelc, s.FuncName.Name(), args,
	)
}

func (s *Call) GenArg() ir.IR {
	return ir.IR(`i32 %%%d`).Expand(s.ResultReg())
}

func (s *Call) GenPrinter() ir.IR {
	return GenIntPrinter(s.Result)
}
