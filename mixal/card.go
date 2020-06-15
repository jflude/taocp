package mixal

import (
	"errors"
	"strings"

	"github.com/jflude/gnuth/mix"
)

var (
	ErrFormatError   = errors.New("format error")
	ErrInternalError = errors.New("internal error")
)

func (a *asmb) processCard(line string, parse parseFunc) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err2, ok := r.(error)
			if !ok || errors.Is(err2, ErrInternalError) {
				panic(r)
			}
			err = a.specifyError(err2, line)
		}
	}()
	a.input = line
	a.count++
	if len(a.input) == 0 {
		panic(ErrFormatError)
	}
	if _, err := mix.ConvertToMIX(a.input); err != nil {
		panic(err)
	}
	if a.input[0] == '*' {
		return nil
	}
	if a.extractColumns(11, 11, true) != "" ||
		a.extractColumns(16, 16, true) != "" {
		panic(ErrFormatError)
	}
	loc := a.extractColumns(1, 10, true)
	op := a.extractColumns(12, 15, true)
	address := a.extractColumns(17, 80, false)
	if op == "ALF" {
		if len(address) > 5 {
			address = address[:5]
		}
	} else if sp := strings.IndexByte(address, ' '); sp != -1 {
		address = address[:sp]
	}
	parse(a, loc, op, address)
	return nil
}
