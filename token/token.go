package token

type Kind int

const (
	KindString = iota
)

type Token struct {
	Kind Kind   // this token kind, it determine what it is.
	Str  string // the substring for this token.
	Beg  int    // the token start position in source code.
	End  int    // the token end position in source code.
}
