package mixal

import "errors"

var ErrUndefinedSymbol = errors.New("undefined symbol")

func (a *asmb) matchAsterisk() bool {
	if len(a.input) > 0 {
		if a.input[0] == '*' {
			a.addToken(asterisk, a.self)
			a.input = a.input[1:]
			return true
		}
	}
	return false
}

func (a *asmb) matchAtomic() bool {
	if a.matchNumber() {
		return true
	}
	if a.matchSymbol() {
		if _, ok := a.symbols[a.lastString()]; !ok {
			a.semanticError(ErrUndefinedSymbol)
		}
		return true
	}
	return a.matchAsterisk()
}

// MIXAL's grammar as described in TAOCP is left-recursive, so these two
// parsers are tweaked to be right-recursive.
func (a *asmb) matchExpr() bool {
	if a.matchAtomic() {
		if a.matchBinaryOp() {
			return a.matchExpr()
		}
		return true
	}
	if a.matchBinaryOp() {
		if k := a.lastKind(); k == '+' || k == '-' {
			return a.matchExpr()
		}
	}
	return false
}

func (a *asmb) matchWValue() bool {
	if a.matchExpr() && a.matchFPart() {
		if a.matchChar(',') {
			a.matchWValue()
		}
		return true
	}
	return false
}
