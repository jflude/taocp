package mixal

import (
	"errors"
	"fmt"
)

var ErrSyntaxError = errors.New("syntax error")

func parseError(err error, text string) {
	panic(fmt.Errorf("%w: %q", err, text))
}

func (a *asmb) syntaxError() {
	parseError(ErrSyntaxError, a.input)
}

func (a *asmb) semanticError(err error) {
	parseError(err, a.lastString())
}

func (a *asmb) specifyError(err error, line string) error {
	return fmt.Errorf("%w in line %d: %s", err, a.count, line)
}
