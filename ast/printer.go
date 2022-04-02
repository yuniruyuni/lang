package ast

import (
	"github.com/yuniruyuni/lang/ir"
)

func GenIntPrinter(reg Reg) ir.IR {
	n := "intfmt"
	l := 4

	return ir.IR(`
		call i32 (i8*, ...) @printf(
			i8* getelementptr inbounds (
				[%d x i8],
				[%d x i8]* @.%s,
				i64 0,
				i64 0
			),
			i32 %%%d
		)
	`).Expand(l, l, n, reg)
}
