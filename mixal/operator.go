package mixal

import "github.com/jflude/gnuth/mix"

func (a *asmb) matchOperator() bool {
	var i int
	for i = 0; i < 4; i++ {
		if i >= len(a.input) ||
			(!mix.IsLetter(rune(a.input[0])) &&
				!mix.IsDigit(rune(a.input[0]))) {
			break
		}
	}
	op := a.input[:i]
	switch op {
	case "EQU", "ORIG", "CON", "ALF", "END":
	default:
		if _, ok := opcodes[op]; !ok {
			return false
		}
	}
	a.addToken(operator, op)
	a.input = a.input[i:]
	return true
}
