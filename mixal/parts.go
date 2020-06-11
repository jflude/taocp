package mixal

func (a *asmb) matchAPart() bool {
	if !a.matchExpr() {
		if a.matchSymbol() {
			_, ok := a.symbols[a.lastString()]
			return !ok
		}
	}
	return true
}

func (a *asmb) matchIPart() bool {
	if a.matchChar(',') {
		return a.matchExpr()
	}
	return true
}

func (a *asmb) matchFPart() bool {
	if a.matchChar('(') {
		return a.matchExpr() && a.matchChar(')')
	}
	return true
}

func (a *asmb) matchLiteral() bool {
	return a.matchChar('=') && a.matchWValue() && a.matchChar('=')
}
