package token

type Kind int

const (
	KindDoubleQuote = iota
	KindString
)

type Token struct {
	Kind Kind
}
