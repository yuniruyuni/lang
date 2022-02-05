package main

import (
	"fmt"
)

const content = `
@.str = private unnamed_addr constant [5 x i8] c"test\00", align 1

define dso_local i32 @main() #0 {
  %1 = alloca i32, align 4
  store i32 0, i32* %1, align 4
  %2 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str, i64 0, i64 0))
  ret i32 0
}

declare i32 @printf(i8*, ...) #1`

func main() {
	fmt.Println(content)
}
