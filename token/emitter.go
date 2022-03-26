package token

import (
	"github.com/yuniruyuni/lang/token/kind"
)

type Emitter func(t *Tokenizer) *Token

func Save(tk *Tokenizer) *Token {
	return nil
}

func Emit(k kind.Kind) Emitter {
	return func(tk *Tokenizer) *Token {
		return &Token{
			Kind: k,
			Str:  tk.code[tk.beg:tk.cur],
			Beg:  tk.beg,
			End:  tk.cur,
		}
	}
}
