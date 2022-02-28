package token

import (
	"github.com/yuniruyuni/lang/token/state"
)

type Tokenizer struct {
	// current state for this tokenizer.
	State state.State

	// entire source code.
	code string

	// the beginning position in source code
	// for current state.
	beg int

	// the current position in source code
	cur int

	// buffer for result tokens.
	tokens []*Token
}

func (t *Tokenizer) emit(tk *Token) {
	t.tokens = append(t.tokens, tk)
}

func (t *Tokenizer) next(pos int, ch rune) {
	t.cur = pos
	table.Run(t, ch)
}

func (t *Tokenizer) Tokenize(code string) []*Token {
	t.code = code
	t.tokens = []*Token{}

	for pos, ch := range t.code {
		t.next(pos, ch)
	}
	// 0 is 0-value of rune, it can assume as NULL char.
	t.next(len(t.code), 0)

	return t.tokens
}
