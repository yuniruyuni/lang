package ir

import (
	"fmt"
	"strings"
)

type IR string

func (ir IR) Expand(vals ...interface{}) IR {
	return IR(fmt.Sprintf(string(ir), vals...))
}

type IRs []IR

// Concat concatinate all given IRs into single IR.
func Concat(irs ...IR) IR {
	strs := make([]string, 0, len(irs))
	for _, ir := range irs {
		strs = append(strs, string(ir))
	}
	return IR(strings.Join(strs, "\n"))
}
