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

// Translate detects reserved words for identifiers and
// it returns new Token that has precisely kind for such identifier.
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
		return t.changeKind(kind.If)
	case "else":
		return t.changeKind(kind.Else)
	case "let":
		return t.changeKind(kind.Let)
	case "while":
		return t.changeKind(kind.While)
	default:
		return t
	}
}

// changeKind returns new token that is changed kind to k.
func (t *Token) changeKind(k kind.Kind) *Token {
	return &Token{
		Kind: k,
		Str:  t.Str,
		Beg:  t.Beg,
		End:  t.End,
	}
}
