package mixal

func (a *asmb) parseAPart() bool {
	if a.exprVal = nil; a.parseExpr() {
		a.aa = *a.exprVal
		return true
	}
	if a.matchSymbol() {
		if _, ok := a.symbols[a.lastString()]; ok {
			return false
		}
		// TODO: mark as a future ref needing a fix-up
	}
	return true
}

func (a *asmb) parseIPart() bool {
	if a.matchChar(',') {
		if a.exprVal = nil; !a.parseExpr() {
			return false
		}
		a.i = a.exprVal.Int()
	}
	return true
}

func (a *asmb) parseFPart() bool {
	if a.matchChar('(') {
		if a.exprVal = nil; !a.parseExpr() || !a.matchChar(')') {
			return false
		}
		a.f = a.exprVal.Int()
	}
	return true
}

func (a *asmb) parseLiteral() bool {
	if !a.matchChar('=') {
		return false
	}
	if a.wVal = 0; !a.parseWValue() || !a.matchChar('=') {
		return false
	}
	// TODO: create future ref for literal and record value to be emitted
	return true
}
