package token

import (
	"github.com/yuniruyuni/lang/token/kind"
	"github.com/yuniruyuni/lang/token/state"
)

type Edges []Edge
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
		{check: Ch(' '), emit: Emit(kind.Skip), next: state.Init, retry: false},
		{check: Ch('\t'), emit: Emit(kind.Skip), next: state.Init, retry: false},
		{check: Ch('"'), emit: Save, next: state.String, retry: false},
		{check: Ch('+'), emit: Emit(kind.Plus), next: state.Init, retry: false},
		{check: Ch('-'), emit: Emit(kind.Minus), next: state.Init, retry: false},
		{check: Ch('*'), emit: Emit(kind.Multiply), next: state.Init, retry: false},
		{check: Ch('/'), emit: Emit(kind.Divide), next: state.Init, retry: false},
		{check: Ch('('), emit: Emit(kind.LeftParen), next: state.Init, retry: false},
		{check: Ch(')'), emit: Emit(kind.RightParen), next: state.Init, retry: false},
		{check: Ch('<'), emit: Emit(kind.Less), next: state.Init, retry: false},
		{check: Ch('='), emit: Emit(kind.Equal), next: state.Init, retry: false},
		{check: IsDigit, emit: Emit(kind.Skip), next: state.Integer, retry: true},
	},
	state.String: Edges{
		{check: Ch('"'), emit: Emit(kind.String), next: state.Init, retry: false},
		{check: Ch('\\'), emit: Save, next: state.Escape, retry: false},
	},
	state.Escape: Edges{
		{check: Any, emit: Save, next: state.String, retry: false},
	},
	state.Integer: Edges{
		{check: IsDigit, emit: Save, next: state.Integer, retry: false},
		{check: Any, emit: Emit(kind.Integer), next: state.Init, retry: true},
	},
}

func (tr Transition) Run(tk *Tokenizer, ch rune) bool {
	return tr[tk.State].Run(tk, ch)
}

func (es Edges) Run(tk *Tokenizer, ch rune) bool {
	for _, e := range es {
		if !e.check(ch) {
			continue
		}

		if !e.retry {
			tk.cur += 1
		}

		t := e.emit(tk)
		if t != nil {
			tk.emit(t)
			tk.beg = tk.cur
		}

		tk.State = e.next

		return e.retry
	}
	return false
}
