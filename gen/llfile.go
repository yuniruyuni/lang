package gen

import (
	"github.com/yuniruyuni/lang/ast"
)

const bodyOpen = `
@.intfmt = private unnamed_addr constant [4 x i8] c"%d\0A\00", align 1

define dso_local i32 @main() #0 {
  %1 = alloca i32, align 4
  store i32 0, i32* %1, align 4
`

const bodyClose = `
   ret i32 0
}

declare i32 @printf(i8*, ...) #1
`

type LLFile struct {
	AST ast.AST
}

func (ll *LLFile) Generate() ast.IR {
	gen := ast.Gen{}
	_ = gen.NextReg()

	ll.AST.AcquireReg(&gen)

	header := ll.AST.GenHeader()
	body := bodyOpen + ll.AST.GenBody() + ll.AST.GenPrinter() + bodyClose

	return ast.IR(header + body)
}
