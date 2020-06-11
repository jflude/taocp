package mixal

import (
	"bufio"
	"io"
	"strings"

	"github.com/jflude/gnuth/mix"
)

func (a *asmb) readProgram(r io.Reader, lp lineParser) error {
	defer func() {
		if err := recover(); err != nil {
			err = a.specifyError(err.(error))
		}
	}()
	a.symbols = make(map[string]int)
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		a.count++
		if len(line) == 0 {
			panic(ErrFormatError)
		}
		if _, err := mix.ConvertToMIX(line); err != nil {
			panic(err)
		}
		if line[0] == '*' {
			continue
		}
		if extractColumns(line, 11, 11, true) != "" ||
			extractColumns(line, 16, 16, true) != "" {
			panic(ErrFormatError)
		}
		loc := extractColumns(line, 1, 10, true)
		op := extractColumns(line, 12, 15, true)
		address := extractColumns(line, 17, 80, false)
		if op == "ALF" {
			if len(address) > 5 {
				address = address[:5]
			}
		} else {
			sp := strings.IndexRune(address, rune(' '))
			if sp != -1 {
				address = address[:sp]
			}
		}
		lp(a, loc, op, address)
	}
	return nil
}
