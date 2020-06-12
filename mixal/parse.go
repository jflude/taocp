package mixal

import "github.com/jflude/gnuth/mix"

func parseLine(a *asmb, loc, op, address string) {
	a.tokens = nil
	a.input = loc
	if a.input != "" {
		if !a.matchSymbol() {
			a.syntaxError()
		}
		sym := a.lastString()
		if isLocalSymbol(sym) {
			if sym[len(sym)-1] != 'H' {
				a.semanticError(ErrInvalidLocal)
			}
		} else if _, ok := a.symbols[sym]; ok {
			a.semanticError(ErrRedefinedSymbol)
		}
	}
	a.input = op
	if !a.matchOperator() {
		a.semanticError(ErrInvalidOperator)
	}
	if a.tokens[0].kind == symbol {
		if a.lastString() != "EQU" {
			a.symbols[a.tokens[0].val.(string)] = a.self
			// TODO: fix-up any future refs seen so far
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
		def := opcodes[a.lastString()]
		a.c, a.f = def.c, def.f
		a.i, a.aa = 0, 0
		if !a.parseAPart() || !a.parseIPart() || !a.parseFPart() {
			a.syntaxError()
		}
		var w mix.Word
		w.PackOp(a.aa, a.i, a.f, a.c)
		a.emit(w)
	}
}
