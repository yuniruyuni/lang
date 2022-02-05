package token

type State int

const (
	StateInit = iota
	StateString
)

type Tokenizer struct {
	// current state for this tokenizer.
	State State

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

func (t *Tokenizer) changeState(st State) {
	t.State = st
	t.beg = t.cur + 1
}

func (t *Tokenizer) forInit(ch rune) {
	if ch == ' ' || ch == '\t' {
		return
	}

	if ch == '"' {
		t.changeState(StateString)
		return
	}
}

func (t *Tokenizer) forString(ch rune) {
	if ch == '"' {
		tk := &Token{
			Kind: KindString,
			Str:  t.code[t.beg:t.cur],
			Beg:  t.beg,
			End:  t.cur,
		}
		t.emit(tk)
		t.changeState(StateInit)
		return
	}
}

func (t *Tokenizer) Tokenize(code string) []*Token {
	t.code = code
	t.tokens = []*Token{}
	for pos, ch := range t.code {
		t.cur = pos
		switch t.State {
		case StateInit:
			t.forInit(ch)
		case StateString:
			t.forString(ch)
		}
	}

	return t.tokens
}
