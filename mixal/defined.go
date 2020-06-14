package mixal

func (a *asmb) matchDefinedSymbol() bool {
	return a.matchSymbol(func(sym *string) bool {
		if isLocalSymbol(*sym) {
			switch (*sym)[1] {
			case 'B':
				*sym = (*sym)[:1] + "H"
				if _, ok := a.symbols[*sym]; ok {
					return true
				}
				fallthrough
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
