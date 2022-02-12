package token

import (
	"github.com/yuniruyuni/lang/token/kind"
)

type Token struct {
	Kind kind.Kind // this token kind, it determine what it is.
	Str  string    // the substring for this token.
	Beg  int       // the token start position in source code.
	End  int       // the token end position in source code.
}
