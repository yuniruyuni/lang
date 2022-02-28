package token

import (
	"github.com/yuniruyuni/lang/token/state"
)

type Edges []Edge
type Transition map[state.State]Edges

type Edge struct {
	check Checker
	emit  Emitter
	next  state.State
}

var table = Transition{
	state.Init: Edges{
		{check: NilCh, emit: Skip, next: state.Init},
		{check: Ch(' '), emit: Skip, next: state.Init},
		{check: Ch('\t'), emit: Skip, next: state.Init},
		{check: Ch('"'), emit: Skip, next: state.String},
		{check: Ch('+'), emit: Skip, next: state.Plus},
		{check: IsDigit, emit: Skip, next: state.Integer},
	},
	state.String: Edges{
		{check: Ch('"'), emit: EmitString, next: state.Init},
		{check: Ch('\\'), emit: Skip, next: state.Escape},
	},
	state.Escape: Edges{
		{check: Any, emit: Skip, next: state.String},
	},
	state.Integer: Edges{
		{check: IsDigit, emit: Skip, next: state.Integer},
		{check: Ch('"'), emit: EmitInteger, next: state.String},
		{check: Ch('+'), emit: EmitInteger, next: state.Plus},
		{check: Any, emit: EmitInteger, next: state.Init},
	},
	state.Plus: Edges{
		{check: Ch('"'), emit: EmitPlus, next: state.String},
		{check: IsDigit, emit: EmitPlus, next: state.Integer},
		{check: Any, emit: EmitPlus, next: state.Init},
	},
}

func (tr Transition) Run(tk *Tokenizer, ch rune) {
	tr[tk.State].Run(tk, ch)
}

func (es Edges) Run(tk *Tokenizer, ch rune) {
	for _, e := range es {
		if !e.check(ch) {
			continue
		}

		t := e.emit(tk)
		if t != nil {
			tk.emit(t)
		}

		if tk.State != e.next {
			tk.beg = tk.cur
		}

		tk.State = e.next

		return
	}
}
