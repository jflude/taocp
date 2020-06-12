package mixal

func (a *asmb) parseAtomic() bool {
	if a.matchNumber() {
		a.evalArg(a.lastQuantity())
		return true
	}
	if a.matchSymbol() {
		if _, ok := a.symbols[a.lastString()]; !ok {
			return false
		}
		a.evalArg(a.lastQuantity())
		return true
	}
	if a.matchAsterisk() {
		a.evalArg(a.lastQuantity())
		return true
	}
	return false
}

// MIXAL's grammar as described in TAOCP is left-recursive and therefore cannot
// be parsed by recursive descent, so parseExpr and parseWValue are modified
// to be right-recursive.
func (a *asmb) parseExpr() bool {
	if a.parseAtomic() {
		if a.matchBinaryOp() {
			a.evalOp(a.lastKind())
			return a.parseExpr()
		}
		return true
	}
	if a.matchBinaryOp() {
		if k := a.lastKind(); k == '+' || k == '-' {
			// convert unary +/- to binary +/- with implied zero
			a.evalArg(0)
			a.evalOp(a.lastKind())
			return a.parseExpr()
		}
	}
	return false
}

func (a *asmb) parseWValue() bool {
	a.f = 5
	if a.exprVal = nil; !a.parseExpr() || !a.parseFPart() {
		return false
	}
	a.wVal.SetField(a.f, *a.exprVal)
	return !a.matchChar(',') || a.parseWValue()
}
