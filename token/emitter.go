package token

import (
	"github.com/yuniruyuni/lang/token/kind"
)

type Emitter func(t *Tokenizer)

func Save(tk *Tokenizer) {}

func Emit(k kind.Kind) Emitter {
	return func(tk *Tokenizer) {
		tk.emit(&Token{
			Kind: k,
			Str:  tk.code[tk.beg:tk.cur],
			Beg:  tk.beg,
			End:  tk.cur,
		})
	}
}
