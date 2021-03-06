package ast

import (
	"strings"

	"github.com/yuniruyuni/lang/ir"
)

type Args struct {
	Result Reg

	// for `x, y, z,`
	Values []AST
}

func (s *Args) Name() Name {
	return ""
}

func (s *Args) Type() Type {
	ts := make([]string, 0, len(s.Values))
	for _, v := range s.Values {
		ts = append(ts, string(v.Type()))
	}
	return Type(strings.Join(ts, ","))
}

func (s *Args) ResultReg() Reg {
	return s.Result
}

func (s *Args) ResultLabel() Label {
	return s.Values[len(s.Values)-1].ResultLabel()
}

func (s *Args) GenHeader(g *Gen) ir.IR {
	hs := make([]ir.IR, 0, len(s.Values))
	for _, v := range s.Values {
		hs = append(hs, v.GenHeader(g))
	}
	return ir.Concat(hs...)
}

func (s *Args) GenBody(g *Gen) ir.IR {
	bodies := make([]ir.IR, 0, len(s.Values))
	for _, v := range s.Values {
		bodies = append(bodies, v.GenBody(g))
	}
	return ir.Concat(bodies...)
}

func (s *Args) GenArg() ir.IR {
	args := make([]ir.IR, 0, len(s.Values))
	for _, v := range s.Values {
		args = append(args, v.GenArg())
	}
	return ir.Join(",", args...)
}

func (s *Args) GenPrinter() ir.IR {
	return GenIntPrinter(s.Result)
}
