package ast

import (
	"strings"

	"github.com/yuniruyuni/lang/ir"
)

type Params struct {
	Result Reg

	// for `x, y, z,`
	Vars []AST
}

func (s *Params) Name() Name {
	return ""
}

func (s *Params) Type() Type {
	ts := make([]string, 0, len(s.Vars))
	for _, v := range s.Vars {
		ts = append(ts, string(v.Type()))
	}
	return Type(strings.Join(ts, ","))
}

func (s *Params) ResultReg() Reg {
	return s.Result
}

func (s *Params) ResultLabel() Label {
	return 0
}

func (s *Params) GenHeader(g *Gen) ir.IR {
	return ""
}

func (s *Params) GenBody(g *Gen) ir.IR {
	bodies := make([]ir.IR, 0, len(s.Vars))
	for _, v := range s.Vars {
		bodies = append(bodies, v.GenBody(g))
	}
	return ir.Join(",", bodies...)
}

func (s *Params) GenArg() ir.IR {
	return ""
}

func (s *Params) GenPrinter() ir.IR {
	return ir.IR("")
}
