package kind

type Kind int

const (
	String = iota
	Integer
	Identifier
	Less
	Equal
	Plus
	Minus
	Multiply
	Divide
	LeftParen
	RightParen
)
