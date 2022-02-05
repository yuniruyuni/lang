package token

type State int

const (
	StateInit = iota
	StateString
)

type Tokenizer struct {
	State
}

func (t *Tokenizer) Tokenize(s string) []*Token {
	return nil
}
