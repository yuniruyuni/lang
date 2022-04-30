package gen

import (
	"github.com/yuniruyuni/lang/ast"
	"github.com/yuniruyuni/lang/ir"
)

const bodyOpen = `
@.intfmt = private unnamed_addr constant [4 x i8] c"%d\0A\00", align 1
@.readfmt = private unnamed_addr constant [4 x i8] c"%d\0A\00", align 1

define i32 @read() {
	%1 = alloca i32, align 4
	store i32 0, i32* %1, align 4
	call i32 (i8*, ...) @scanf(
		i8* getelementptr inbounds (
			[4 x i8],
			[4 x i8]* @.readfmt,
			i64 0,
			i64 0
		),
		i32* %1
	)
	%3 = load i32, i32* %1, align 4
	ret i32 %3
}

declare i32 @scanf(i8*, ...)
declare i32 @printf(i8*, ...)

define i32 @main() {
	%1 = alloca i32, align 4
	store i32 0, i32* %1, align 4
`

const bodyClose = `
	ret i32 0
}
`

type LLFile struct {
	AST ast.AST
}

func (ll *LLFile) Generate() ir.IR {
	gen := ast.NewGen()
	_ = gen.NextReg()

	gen.RegisterFunc("printf", ast.Type{"i8*", "..."})
	gen.RegisterFunc("read", ast.Type{"i8*", "..."})

	header := ll.AST.GenHeader(gen)
	body := bodyOpen + ll.AST.GenBody(gen) + ll.AST.GenPrinter() + bodyClose

	return ir.IR(header + body)
}
