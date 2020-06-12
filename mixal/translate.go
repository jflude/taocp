package mixal

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/jflude/gnuth/mix"
)

func (a *asmb) translate(r io.Reader,
	parser func(*asmb, string, string, string)) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if err = r.(error); errors.Is(err, ErrInternalError) {
				panic(err)
			}
			err = a.specifyError(err)
		}
	}()
	a.symbols = make(map[string]int)
	a.newSegment(0)
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		a.input = sc.Text()
		a.count++
		if len(a.input) == 0 {
			panic(ErrFormatError)
		}
		if _, err := mix.ConvertToMIX(a.input); err != nil {
			panic(err)
		}
		if a.input[0] == '*' {
			continue
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
		parser(a, loc, op, address)
	}
	return nil
}
