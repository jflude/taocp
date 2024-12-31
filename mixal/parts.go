package mixal

import "fmt"

func (a *asmb) parseAPart() bool {
	if a.parseLiteral() {
		return true
	}
	if a.matchFutureRef() {
		a.addFixUp(a.lastString())
		return true
	}
	if a.exprVal = nil; a.parseExpr() {
		a.aa = *a.exprVal
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
		save := a.exprVal
		if a.exprVal = nil; !a.parseExpr() || !a.matchChar(')') {
			return false
		}
		a.f = a.exprVal.Int()
		a.exprVal = save
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
	sym := fmt.Sprintf("_%d", a.label)
	a.label++
	a.literals = append(a.literals, literal{sym, a.wVal})
	a.addFixUp(sym)
	return true
}
