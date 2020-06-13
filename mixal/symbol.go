package mixal

import "github.com/jflude/gnuth/mix"

func (a *asmb) matchDefinedSymbol() bool {
	return a.matchSymbol(true)
}

func (a *asmb) matchFutureRef() bool {
	return a.matchSymbol(false)
}

func (a *asmb) matchSymbol(defined bool) bool {
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
	if _, ok := a.symbols[a.input[:i]]; ok != defined {
		return false
	}
	a.addToken(symbol, a.input[:i])
	a.input = a.input[i:]
	return true
}
