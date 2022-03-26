package token

import "unicode"

type Checker func(ch rune) bool

func Ch(want rune) Checker {
	return func(ch rune) bool { return ch == want }
}

func NilCh(ch rune) bool {
	return ch == 0
}

func Any(ch rune) bool {
	return true
}

func IsDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func IsLetter(ch rune) bool {
	return unicode.IsLetter(ch)
}
