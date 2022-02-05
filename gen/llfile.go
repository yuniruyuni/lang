package gen

import (
	"bufio"
	"fmt"
	"os"
)

const content = `
@.str = private unnamed_addr constant [%d x i8] c"%s\00", align 1

define dso_local i32 @main() #0 {
  %%1 = alloca i32, align 4
  store i32 0, i32* %%1, align 4
  %%2 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([%d x i8], [%d x i8]* @.str, i64 0, i64 0))
  ret i32 0
}

declare i32 @printf(i8*, ...) #1`

var sc = bufio.NewScanner(os.Stdin)

type LLFile struct {
	Word string
}

const (
	nullCharSize = 1
)

func (ll *LLFile) Generate() string {
	// it needs to calculate byte size.
	l := len(ll.Word) + nullCharSize
	w := ll.Word
	return fmt.Sprintf(content, l, w, l, l)
}
