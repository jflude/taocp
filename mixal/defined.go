package mixal

func (a *asmb) matchDefinedSymbol() bool {
	return a.matchSymbol(func(sym *string) bool {
		if isLocalSymbol(*sym) {
			switch (*sym)[1] {
			case 'B':
				*sym = (*sym)[:1] + "H"
			case 'F':
				return false
			case 'H':
				parseError(ErrInvalidLocal, *sym)
			}
		}
		_, ok := a.symbols[*sym]
		return ok
	})
}
