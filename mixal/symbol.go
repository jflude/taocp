package mixal

import "github.com/jflude/gnuth/mix"

func (a *asmb) matchSymbol() bool {
	var letter bool
	var i int
	for i = 0; i < 10; i++ {
		if i >= len(a.input) {
			break
		}
		if mix.IsLetter(rune(a.input[i])) {
			letter = true
		} else if !mix.IsDigit(rune(a.input[i])) {
			break
		}
	}
	if !letter {
		return false
	}
	a.addToken(symbol, a.input[:i])
	a.input = a.input[i:]
	return true
}

func isLocalSymbol(sym string) bool {
	return len(sym) == 2 && mix.IsDigit(rune(sym[0])) &&
		(sym[1] == 'B' || sym[1] == 'F' || sym[1] == 'H')
}
