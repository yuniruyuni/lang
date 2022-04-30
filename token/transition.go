package token

import (
	"github.com/yuniruyuni/lang/token/kind"
	"github.com/yuniruyuni/lang/token/state"
)

type Edges []*Edge
type Transition map[state.State]Edges

type Edge struct {
	check Checker
	emit  Emitter
	next  state.State
	retry bool // if true, recheck same character on next state.
}

var table = Transition{
	state.Init: Edges{
		{check: NilCh, emit: Save, next: state.Init, retry: false},
		{check: Ch('"'), emit: Save, next: state.String, retry: false},
		{check: Ch('+'), emit: Emit(kind.Plus), next: state.Init, retry: false},
		{check: Ch('-'), emit: Emit(kind.Minus), next: state.Init, retry: false},
		{check: Ch('*'), emit: Emit(kind.Multiply), next: state.Init, retry: false},
		{check: Ch('/'), emit: Emit(kind.Divide), next: state.Init, retry: false},
		{check: Ch('('), emit: Emit(kind.LeftParen), next: state.Init, retry: false},
		{check: Ch(')'), emit: Emit(kind.RightParen), next: state.Init, retry: false},
		{check: Ch('{'), emit: Emit(kind.LeftCurly), next: state.Init, retry: false},
		{check: Ch('}'), emit: Emit(kind.RightCurly), next: state.Init, retry: false},
		{check: Ch('<'), emit: Emit(kind.Less), next: state.Init, retry: false},
		{check: Ch('='), emit: Emit(kind.Equal), next: state.Init, retry: false},
		{check: Ch(';'), emit: Emit(kind.Semicolon), next: state.Init, retry: false},
		{check: Ch(','), emit: Emit(kind.Comma), next: state.Init, retry: false},
		{check: IsDigit, emit: Save, next: state.Integer, retry: true},
		{check: IsLetter, emit: Save, next: state.Identifier, retry: true},
		{check: Any, emit: Emit(kind.Skip), next: state.Init, retry: false},
	},
	state.String: Edges{
		{check: Ch('"'), emit: Emit(kind.String), next: state.Init, retry: false},
		{check: Ch('\\'), emit: Save, next: state.Escape, retry: false},
		{check: Any, emit: Save, next: state.String, retry: false},
	},
	state.Escape: Edges{
		{check: Any, emit: Save, next: state.String, retry: false},
	},
	state.Integer: Edges{
		{check: IsDigit, emit: Save, next: state.Integer, retry: false},
		{check: Any, emit: Emit(kind.Integer), next: state.Init, retry: true},
	},
	state.Identifier: Edges{
		{check: IsLetter, emit: Save, next: state.Identifier, retry: false},
		{check: Any, emit: Emit(kind.Identifier), next: state.Init, retry: true},
	},
}

func (tr Transition) Run(tk *Tokenizer, ch rune) bool {
	return tr[tk.State].Run(tk, ch)
}

func (es Edges) edgeFor(ch rune) *Edge {
	for _, e := range es {
		if e.check(ch) {
			return e
		}
	}
	return nil
}

func (es Edges) Run(tk *Tokenizer, ch rune) bool {
	e := es.edgeFor(ch)
	if e == nil {
		panic("There is no edge for next token")
	}

	if !e.retry {
		tk.cur += 1
	}

	e.emit(tk)

	tk.State = e.next

	return e.retry
}
