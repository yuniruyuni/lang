package state

type State int

const (
	Init = iota
	String
	Escape
	Integer
)
