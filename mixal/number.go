package mixal

import (
	"strconv"

	"github.com/jflude/gnuth/mix"
)

func (a *asmb) matchNumber() bool {
	var i int
	for i = 0; i < 10; i++ {
		if i >= len(a.input) || !mix.IsDigit(rune(a.input[i])) {
			break
		}
	}
	if i == 0 {
		return false
	}
	n, err := strconv.Atoi(a.input[:i])
	if err != nil {
		panic(err)
	}
	a.addToken(number, n)
	a.input = a.input[i:]
	return true
}
