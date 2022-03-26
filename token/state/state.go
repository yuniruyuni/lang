package state

type State int

const (
	Init = iota
	String
	Escape
	Less
	Equal
	Integer
	Plus
	Minus
	Multiply
	Divide
	LeftParen
	RightParen
)
