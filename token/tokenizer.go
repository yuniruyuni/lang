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
	t.beg = t.cur
}

func IsDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func (t *Tokenizer) forInit(ch rune) {
	if ch == ' ' || ch == '\t' || ch == 0 {
		return
	}

	if ch == '"' {
		t.changeState(state.String)
		return
	}

	if ch == '+' {
		t.changeState(state.Plus)
		return
	}

	if IsDigit(ch) {
		t.changeState(state.Integer)
		return
	}

	// TODO: emit error
}

func (t *Tokenizer) forString(ch rune) {
	if ch == '"' {
		tk := &Token{
			Kind: kind.String,
			Str:  t.code[t.beg+1 : t.cur],
			Beg:  t.beg + 1,
			End:  t.cur,
		}
		t.emit(tk)
		t.changeState(state.Init)
		return
	}

	// read all next token
}

func (t *Tokenizer) forInteger(ch rune) {
	if IsDigit(ch) {
		// if 0~9(digit char) has come,
		// we will continue to read it to read next digit.
		return
	}

	// for other char, we assume it is the end of an integer constant.
	tk := &Token{
		Kind: kind.Integer,
		Str:  t.code[t.beg:t.cur],
		Beg:  t.beg,
		End:  t.cur,
	}
	t.emit(tk)

	if ch == '"' {
		t.changeState(state.String)
		return
	}

	if ch == '+' {
		t.changeState(state.Plus)
		return
	}

	t.changeState(state.Init)
}

func (t *Tokenizer) forPlus(ch rune) {
	tk := &Token{
		Kind: kind.Plus,
		Str:  t.code[t.beg:t.cur],
		Beg:  t.beg,
		End:  t.cur,
	}
	t.emit(tk)

	if ch == '"' {
		t.changeState(state.String)
		return
	}

	if IsDigit(ch) {
		t.changeState(state.Integer)
		return
	}

	t.changeState(state.Init)
}

func (t *Tokenizer) next(pos int, ch rune) {
	t.cur = pos
	switch t.State {
	case state.Init:
		t.forInit(ch)
	case state.String:
		t.forString(ch)
	case state.Integer:
		t.forInteger(ch)
	case state.Plus:
		t.forPlus(ch)
	}
}

func (t *Tokenizer) Tokenize(code string) []*Token {
	t.code = code
	t.tokens = []*Token{}

	var pos int
	var ch rune
	for pos, ch = range t.code {
		t.next(pos, ch)
	}
	// 0 is 0-value of rune, it can assume as NULL char.
	t.next(pos+1, 0)

	return t.tokens
}
