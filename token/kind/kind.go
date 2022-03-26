package kind

type Kind int

const (
	// -------- Virtual Tokens
	Skip = iota

	// -------- Concrete Tokens
	String
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
	LeftCurly
	RightCurly
	If
	Else
)
