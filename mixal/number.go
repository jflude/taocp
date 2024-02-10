package mixal

import (
	"strconv"

	"github.com/jflude/taocp/mix"
)

func (a *asmb) matchNumber() bool {
	var i int
	var r rune
	for i = 0; i < 10 && i < len(a.input); i++ {
		if r = rune(a.input[i]); !mix.IsDigit(r) {
			break
		}
	}
	if i == 0 || mix.IsLetter(r) {
		return false
	}
	n, err := strconv.Atoi(a.input[:i])
	if err != nil {
		panic(ErrInternal)
	}
	a.addToken(number, n)
	a.input = a.input[i:]
	return true
}
