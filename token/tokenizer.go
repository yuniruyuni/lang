package token

import (
	"github.com/yuniruyuni/lang/token/kind"
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

func (t *Tokenizer) changeState(st state.State) {
	t.State = st
	t.beg = t.cur + 1
}

func (t *Tokenizer) forInit(ch rune) {
	if ch == ' ' || ch == '\t' {
		return
	}

	if ch == '"' {
		t.changeState(state.String)
		return
	}
}

func (t *Tokenizer) forString(ch rune) {
	if ch == '"' {
		tk := &Token{
			Kind: kind.String,
			Str:  t.code[t.beg:t.cur],
			Beg:  t.beg,
			End:  t.cur,
		}
		t.emit(tk)
		t.changeState(state.Init)
		return
	}
}

func (t *Tokenizer) Tokenize(code string) []*Token {
	t.code = code
	t.tokens = []*Token{}
	for pos, ch := range t.code {
		t.cur = pos
		switch t.State {
		case state.Init:
			t.forInit(ch)
		case state.String:
			t.forString(ch)
		}
	}

	return t.tokens
}
