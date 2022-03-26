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
		{check: Ch('-'), emit: Skip, next: state.Minus},
		{check: Ch('*'), emit: Skip, next: state.Multiply},
		{check: Ch('/'), emit: Skip, next: state.Divide},
		{check: Ch('('), emit: Skip, next: state.LeftParen},
		{check: Ch(')'), emit: Skip, next: state.RightParen},
		{check: Ch('<'), emit: Skip, next: state.Less},
		{check: Ch('='), emit: Skip, next: state.Equal},
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
		{check: Ch('-'), emit: EmitInteger, next: state.Minus},
		{check: Ch('*'), emit: EmitInteger, next: state.Multiply},
		{check: Ch('/'), emit: EmitInteger, next: state.Divide},
		{check: Ch('('), emit: EmitInteger, next: state.LeftParen},
		{check: Ch(')'), emit: EmitInteger, next: state.RightParen},
		{check: Ch('<'), emit: EmitInteger, next: state.Less},
		{check: Ch('='), emit: EmitInteger, next: state.Equal},
		{check: Any, emit: EmitInteger, next: state.Init},
	},
	state.Plus: Edges{
		{check: Ch('"'), emit: EmitPlus, next: state.String},
		{check: Ch('+'), emit: EmitPlus, next: state.Plus},
		{check: Ch('-'), emit: EmitPlus, next: state.Minus},
		{check: Ch('*'), emit: EmitPlus, next: state.Multiply},
		{check: Ch('/'), emit: EmitPlus, next: state.Divide},
		{check: Ch('('), emit: EmitPlus, next: state.LeftParen},
		{check: Ch(')'), emit: EmitPlus, next: state.RightParen},
		{check: Ch('<'), emit: EmitPlus, next: state.Less},
		{check: Ch('='), emit: EmitPlus, next: state.Equal},
		{check: IsDigit, emit: EmitPlus, next: state.Integer},
		{check: Any, emit: EmitPlus, next: state.Init},
	},
	state.Minus: Edges{
		{check: Ch('"'), emit: EmitMinus, next: state.String},
		{check: Ch('+'), emit: EmitMinus, next: state.Plus},
		{check: Ch('-'), emit: EmitMinus, next: state.Minus},
		{check: Ch('*'), emit: EmitMinus, next: state.Multiply},
		{check: Ch('/'), emit: EmitMinus, next: state.Divide},
		{check: Ch('('), emit: EmitMinus, next: state.LeftParen},
		{check: Ch(')'), emit: EmitMinus, next: state.RightParen},
		{check: Ch('<'), emit: EmitMinus, next: state.Less},
		{check: Ch('='), emit: EmitMinus, next: state.Equal},
		{check: IsDigit, emit: EmitMinus, next: state.Integer},
		{check: Any, emit: EmitMinus, next: state.Init},
	},
	state.Multiply: Edges{
		{check: Ch('"'), emit: EmitMultiply, next: state.String},
		{check: Ch('+'), emit: EmitMultiply, next: state.Plus},
		{check: Ch('-'), emit: EmitMultiply, next: state.Minus},
		{check: Ch('*'), emit: EmitMultiply, next: state.Multiply},
		{check: Ch('/'), emit: EmitMultiply, next: state.Divide},
		{check: Ch('('), emit: EmitMultiply, next: state.LeftParen},
		{check: Ch(')'), emit: EmitMultiply, next: state.RightParen},
		{check: Ch('<'), emit: EmitMultiply, next: state.Less},
		{check: Ch('='), emit: EmitMultiply, next: state.Equal},
		{check: IsDigit, emit: EmitMultiply, next: state.Integer},
		{check: Any, emit: EmitMultiply, next: state.Init},
	},
	state.Divide: Edges{
		{check: Ch('"'), emit: EmitDivide, next: state.String},
		{check: Ch('+'), emit: EmitDivide, next: state.Plus},
		{check: Ch('-'), emit: EmitDivide, next: state.Minus},
		{check: Ch('*'), emit: EmitDivide, next: state.Multiply},
		{check: Ch('/'), emit: EmitDivide, next: state.Divide},
		{check: Ch('('), emit: EmitDivide, next: state.LeftParen},
		{check: Ch(')'), emit: EmitDivide, next: state.RightParen},
		{check: Ch('<'), emit: EmitDivide, next: state.Less},
		{check: Ch('='), emit: EmitDivide, next: state.Equal},
		{check: IsDigit, emit: EmitDivide, next: state.Integer},
		{check: Any, emit: EmitDivide, next: state.Init},
	},
	state.LeftParen: Edges{
		{check: Ch('"'), emit: EmitLeftParen, next: state.String},
		{check: Ch('+'), emit: EmitLeftParen, next: state.Plus},
		{check: Ch('-'), emit: EmitLeftParen, next: state.Minus},
		{check: Ch('*'), emit: EmitLeftParen, next: state.Multiply},
		{check: Ch('/'), emit: EmitLeftParen, next: state.Divide},
		{check: Ch('('), emit: EmitLeftParen, next: state.LeftParen},
		{check: Ch(')'), emit: EmitLeftParen, next: state.RightParen},
		{check: Ch('<'), emit: EmitLeftParen, next: state.Less},
		{check: Ch('='), emit: EmitLeftParen, next: state.Equal},
		{check: IsDigit, emit: EmitLeftParen, next: state.Integer},
		{check: Any, emit: EmitLeftParen, next: state.Init},
	},
	state.RightParen: Edges{
		{check: Ch('"'), emit: EmitRightParen, next: state.String},
		{check: Ch('+'), emit: EmitRightParen, next: state.Plus},
		{check: Ch('-'), emit: EmitRightParen, next: state.Minus},
		{check: Ch('*'), emit: EmitRightParen, next: state.Multiply},
		{check: Ch('/'), emit: EmitRightParen, next: state.Divide},
		{check: Ch('('), emit: EmitRightParen, next: state.LeftParen},
		{check: Ch(')'), emit: EmitRightParen, next: state.RightParen},
		{check: Ch('<'), emit: EmitRightParen, next: state.Less},
		{check: Ch('='), emit: EmitRightParen, next: state.Equal},
		{check: IsDigit, emit: EmitRightParen, next: state.Integer},
		{check: Any, emit: EmitRightParen, next: state.Init},
	},
	state.Less: Edges{
		{check: Ch('"'), emit: EmitLess, next: state.String},
		{check: Ch('+'), emit: EmitLess, next: state.Plus},
		{check: Ch('-'), emit: EmitLess, next: state.Minus},
		{check: Ch('*'), emit: EmitLess, next: state.Multiply},
		{check: Ch('/'), emit: EmitLess, next: state.Divide},
		{check: Ch('('), emit: EmitLess, next: state.LeftParen},
		{check: Ch(')'), emit: EmitLess, next: state.RightParen},
		{check: Ch('<'), emit: EmitLess, next: state.Less},
		{check: Ch('='), emit: EmitLess, next: state.Equal},
		{check: IsDigit, emit: EmitLess, next: state.Integer},
		{check: Any, emit: EmitLess, next: state.Init},
	},
	state.Equal: Edges{
		{check: Ch('"'), emit: EmitEqual, next: state.String},
		{check: Ch('+'), emit: EmitEqual, next: state.Plus},
		{check: Ch('-'), emit: EmitEqual, next: state.Minus},
		{check: Ch('*'), emit: EmitEqual, next: state.Multiply},
		{check: Ch('/'), emit: EmitEqual, next: state.Divide},
		{check: Ch('('), emit: EmitEqual, next: state.LeftParen},
		{check: Ch(')'), emit: EmitEqual, next: state.RightParen},
		{check: Ch('<'), emit: EmitEqual, next: state.Less},
		{check: Ch('='), emit: EmitEqual, next: state.Equal},
		{check: IsDigit, emit: EmitEqual, next: state.Integer},
		{check: Any, emit: EmitEqual, next: state.Init},
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
			tk.beg = tk.cur
		}

		if tk.State != e.next {
			tk.beg = tk.cur
		}

		tk.State = e.next

		return
	}
}
