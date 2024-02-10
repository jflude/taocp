package mixal

import "errors"

var (
	ErrInvalidLocal    = errors.New("mixal: invalid local symbol")
	ErrRedefinedSymbol = errors.New("mixal: redefined symbol")
)

func (a *asmb) matchUndefinedSymbol() bool {
	return a.matchSymbol(func(sym *string) bool {
		if isLocalSymbol(*sym) {
			if (*sym)[1] != 'H' {
				parseError(ErrInvalidLocal, *sym)
			}
		} else if _, ok := a.symbols[*sym]; ok {
			parseError(ErrRedefinedSymbol, *sym)
		}
		return true
	})
}
