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

func IsDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

type Emitter func(t *Tokenizer) *Token

func EmitNil(_ *Tokenizer) *Token {
	return nil
}

type Checker func(ch rune) bool

func Ch(want rune) Checker {
	return func(ch rune) bool { return ch == want }
}

func NilCh(ch rune) bool {
	return ch == 0
}

func Any(ch rune) bool {
	return true
}

type Edge struct {
	check   Checker
	emit    Emitter
	forward bool
	next    state.State
}

type Transition []Edge
type StateMap map[state.State]Transition

var tr = StateMap{
	state.Init: Transition{
		{check: NilCh, emit: EmitNil, next: state.Init},
		{check: Ch(' '), emit: EmitNil, next: state.Init},
		{check: Ch('\t'), emit: EmitNil, next: state.Init},
		{check: Ch('"'), emit: EmitNil, forward: true, next: state.String},
		{check: Ch('+'), emit: EmitNil, forward: true, next: state.Plus},
		{check: IsDigit, emit: EmitNil, forward: true, next: state.Integer},
	},
	state.String: Transition{
		{check: Ch('"'), emit: EmitString, forward: true, next: state.Init},
	},
	state.Integer: Transition{
		{check: IsDigit, emit: EmitNil, next: state.Integer},
		{check: Ch('"'), emit: EmitInteger, forward: true, next: state.String},
		{check: Ch('+'), emit: EmitInteger, forward: true, next: state.Plus},
		{check: Any, emit: EmitInteger, forward: true, next: state.Init},
	},
	state.Plus: Transition{
		{check: Ch('"'), emit: EmitPlus, forward: true, next: state.String},
		{check: IsDigit, emit: EmitPlus, forward: true, next: state.Integer},
		{check: Any, emit: EmitPlus, forward: true, next: state.Init},
	},
}

func (tr Transition) Run(tk *Tokenizer, ch rune) {
	for _, edge := range tr {
		if !edge.check(ch) {
			continue
		}

		t := edge.emit(tk)
		if t != nil {
			tk.emit(t)
		}

		if edge.forward {
			tk.beg = tk.cur
		}
		tk.State = edge.next

		return
	}
}

func EmitString(tk *Tokenizer) *Token {
	return &Token{
		Kind: kind.String,
		Str:  tk.code[tk.beg+1 : tk.cur],
		Beg:  tk.beg + 1,
		End:  tk.cur,
	}
}

func EmitInteger(tk *Tokenizer) *Token {
	return &Token{
		Kind: kind.Integer,
		Str:  tk.code[tk.beg:tk.cur],
		Beg:  tk.beg,
		End:  tk.cur,
	}
}

func EmitPlus(tk *Tokenizer) *Token {
	return &Token{
		Kind: kind.Plus,
		Str:  tk.code[tk.beg:tk.cur],
		Beg:  tk.beg,
		End:  tk.cur,
	}
}

func (t *Tokenizer) next(pos int, ch rune) {
	t.cur = pos
	tr[t.State].Run(t, ch)
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
