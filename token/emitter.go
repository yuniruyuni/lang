package token

import (
	"github.com/yuniruyuni/lang/token/kind"
)

type Emitter func(t *Tokenizer) *Token

func Skip(_ *Tokenizer) *Token {
	return nil
}

func EmitString(tk *Tokenizer) *Token {
	t := &Token{
		Kind: kind.String,
		Str:  tk.code[tk.beg+1 : tk.cur],
		Beg:  tk.beg + 1,
		End:  tk.cur,
	}
	return t
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

func EmitMinus(tk *Tokenizer) *Token {
	return &Token{
		Kind: kind.Minus,
		Str:  tk.code[tk.beg:tk.cur],
		Beg:  tk.beg,
		End:  tk.cur,
	}
}

func EmitMultiply(tk *Tokenizer) *Token {
	return &Token{
		Kind: kind.Multiply,
		Str:  tk.code[tk.beg:tk.cur],
		Beg:  tk.beg,
		End:  tk.cur,
	}
}

func EmitDivide(tk *Tokenizer) *Token {
	return &Token{
		Kind: kind.Divide,
		Str:  tk.code[tk.beg:tk.cur],
		Beg:  tk.beg,
		End:  tk.cur,
	}
}
