package ast

type Integer struct {
	Result Reg
	Label  Label
	Alloc  Reg
	Value  int
}

func (nd *Integer) ResultReg() Reg {
	return nd.Result
}

func (s *Integer) ResultLabel() Label {
	return s.Label
}

func (nd *Integer) GenHeader() IR {
	return "\n"
}

func (nd *Integer) GenBody(g *Gen) IR {
	nd.Alloc = g.NextReg()
	nd.Result = g.NextReg()
	nd.Label = g.CurLabel()

	return IR(`
		%%%d = alloca i32, align 4
		store i32 %d, i32* %%%d
		%%%d = load i32, i32* %%%d, align 4
	`).Expand(
		nd.Alloc, nd.Value, nd.Alloc, nd.Result, nd.Alloc,
	)
}

func (nd *Integer) GenPrinter() IR {
	n := "intfmt"
	l := 4
	v := nd.Result

	return IR(`call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([%d x i8], [%d x i8]* @.%s, i64 0, i64 0), i32 %%%d)`).
		Expand(l, l, n, v)
}
