package mixal

import "errors"

var ErrUndefinedLocal = errors.New("undefined local symbol")

func (a *asmb) matchDefinedSymbol() bool {
	return a.matchSymbol(func(sym *string) bool {
		if isLocalSymbol(*sym) {
			switch (*sym)[1] {
			case 'B':
				save := *sym
				*sym = (*sym)[:1] + "H"
				if _, ok := a.symbols[*sym]; ok {
					return true
				}
				parseError(ErrUndefinedLocal, save)
			case 'H':
				parseError(ErrInvalidLocal, *sym)
			case 'F':
				return false
			}
		}
		_, ok := a.symbols[*sym]
		return ok
	})
}
