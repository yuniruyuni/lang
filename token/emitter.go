package token

import (
	"github.com/yuniruyuni/lang/token/kind"
)

type Emitter func(t *Tokenizer) *Token

func Save(tk *Tokenizer) *Token {
	return nil
}

func Skip(tk *Tokenizer) *Token {
	return &Token{
		Kind: kind.Skip,
		Str:  tk.code[tk.beg:tk.cur],
		Beg:  tk.beg,
		End:  tk.cur,
	}
}

func EmitString(tk *Tokenizer) *Token {
	t := &Token{
		Kind: kind.String,
		Str:  tk.code[tk.beg+1 : tk.cur-1],
		Beg:  tk.beg + 1,
		End:  tk.cur - 1,
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

func EmitLeftParen(tk *Tokenizer) *Token {
	return &Token{
		Kind: kind.LeftParen,
		Str:  tk.code[tk.beg:tk.cur],
		Beg:  tk.beg,
		End:  tk.cur,
	}
}

func EmitRightParen(tk *Tokenizer) *Token {
	return &Token{
		Kind: kind.RightParen,
		Str:  tk.code[tk.beg:tk.cur],
		Beg:  tk.beg,
		End:  tk.cur,
	}
}

func EmitLess(tk *Tokenizer) *Token {
	return &Token{
		Kind: kind.Less,
		Str:  tk.code[tk.beg:tk.cur],
		Beg:  tk.beg,
		End:  tk.cur,
	}
}

func EmitEqual(tk *Tokenizer) *Token {
	return &Token{
		Kind: kind.Equal,
		Str:  tk.code[tk.beg:tk.cur],
		Beg:  tk.beg,
		End:  tk.cur,
	}
}
