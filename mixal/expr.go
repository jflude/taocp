package mixal

func (a *asmb) parseAtomic() bool {
	if a.matchNumber() || a.matchDefinedSymbol() || a.matchAsterisk() {
		a.evalArg(a.lastQuantity())
		return true
	}
	return false
}

// MIXAL's grammar as described in TAOCP is left-recursive and therefore
// cannot be easily parsed top-down, so refactor it to remove the recursion.
func (a *asmb) parseExprDash() bool {
	if a.matchBinaryOp() {
		a.exprOp = a.lastKind()
		return a.parseExpr() && a.parseExprDash()
	}
	return true
}

func (a *asmb) parseExpr() bool {
	if a.parseAtomic() {
		return a.parseExprDash()
	}
	if a.matchUnaryOp() {
		if k := a.lastKind(); a.parseAtomic() {
			if k == '-' {
				*a.exprVal = a.exprVal.Negate()
			}
			return a.parseExprDash()
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
