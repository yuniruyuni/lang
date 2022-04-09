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

func (t *Token) Translate() *Token {
	if t == nil {
		return nil
	}

	if t.Kind == kind.Skip {
		return nil
	}

	if t.Kind != kind.Identifier {
		return t
	}

	switch t.Str {
	case "if":
		return &Token{
			Kind: kind.If,
			Str:  t.Str,
			Beg:  t.Beg,
			End:  t.End,
		}
	case "else":
		return &Token{
			Kind: kind.Else,
			Str:  t.Str,
			Beg:  t.Beg,
			End:  t.End,
		}
	default:
		return t
	}
}
