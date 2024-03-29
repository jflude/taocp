package mixal

import (
	"errors"
	"fmt"
)

var ErrSyntax = errors.New("mixal: syntax error")

func parseError(err error, text string) {
	panic(fmt.Errorf("%w: %q", err, text))
}

func (a *asmb) syntaxError() {
	parseError(ErrSyntax, a.input)
}

func (a *asmb) semanticError(err error) {
	parseError(err, a.lastString())
}

func (a *asmb) specifyError(err error, line string) error {
	return fmt.Errorf("%w in line %d: %s", err, a.count, line)
}
