package mixal

import (
	"errors"
	"fmt"
)

var (
	ErrFormatError     = errors.New("format error")
	ErrSyntaxError     = errors.New("syntax error")
	ErrRedefinedSymbol = errors.New("redefined symbol")
	ErrInvalidLocal    = errors.New("invalid local symbol")
	ErrInternalError   = errors.New("internal error")
)

func parseError(err error, text string) {
	panic(fmt.Errorf("%w: %q", err, text))
}

func (a *asmb) syntaxError() {
	parseError(ErrSyntaxError, a.input)
}

func (a *asmb) semanticError(err error) {
	parseError(err, a.lastString())
}

func (a *asmb) specifyError(err error) error {
	return fmt.Errorf("%w in line %d", err, a.count)
}
