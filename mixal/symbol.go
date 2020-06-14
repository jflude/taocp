package mixal

import "github.com/jflude/gnuth/mix"

func (a *asmb) matchSymbol(criterion func(*string) bool) bool {
	var letter bool
	var i int
	for i = 0; i < 10 && i < len(a.input); i++ {
		if mix.IsLetter(rune(a.input[i])) {
			letter = true
		} else if !mix.IsDigit(rune(a.input[i])) {
			break
		}
	}
	if !letter {
		return false
	}
	sym := a.input[:i]
	if !criterion(&sym) {
		return false
	}
	a.addToken(symbol, sym)
	a.input = a.input[i:]
	return true
}

func isLocalSymbol(sym string) bool {
	if len(sym) == 2 && mix.IsDigit(rune(sym[0])) {
		switch sym[1] {
		case 'B', 'F', 'H':
			return true
		}
	}
	return false
}
