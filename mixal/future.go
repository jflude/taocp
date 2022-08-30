// Copyright (c) 2020 Justin Flude. All rights reserved.
// Use of this source code is governed by the COPYING.md file.
package mixal

func (a *asmb) matchFutureRef() bool {
	return a.matchSymbol(func(sym *string) bool {
		if isLocalSymbol(*sym) {
			switch (*sym)[1] {
			case 'B':
				return false
			case 'F':
				*sym = (*sym)[:1] + "H"
			case 'H':
				parseError(ErrInvalidLocal, *sym)
			}
		} else if _, ok := a.symbols[*sym]; ok {
			return false
		}
		return true
	})
}
