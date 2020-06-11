package mixal

func parseLine(a *asmb, loc, op, address string) {
	a.tokens = nil
	a.input = loc
	if a.input != "" {
		if !a.matchSymbol() {
			a.syntaxError()
		}
		if _, ok := a.symbols[a.lastString()]; ok {
			a.semanticError(ErrRedefinedSymbol)
		}
	}
	a.input = op
	if !a.matchOperator() {
		a.semanticError(ErrInvalidOperator)
	}
	if a.tokens[0].kind == symbol {
		if s := a.lastString(); s != "EQU" {
			a.symbols[s] = a.self
		}
	}
	a.input = address
	switch a.lastString() {
	case "EQU":
		a.parseEQU()
	case "ORIG":
		a.parseORIG()
	case "CON":
		a.parseCON()
	case "ALF":
		a.parseALF()
	case "END":
		a.parseEND()
	default:
		if !a.matchAPart() || !a.matchIPart() || !a.matchFPart() {
			a.syntaxError()
		}
	}
}
