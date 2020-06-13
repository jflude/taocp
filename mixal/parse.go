package mixal

import "github.com/jflude/gnuth/mix"

func parseLine(a *asmb, loc, op, address string) {
	a.tokens = nil
	a.input = loc
	if a.input != "" && !a.matchUndefinedSymbol() {
		a.syntaxError()
	}
	a.input = op
	if !a.matchOperator() {
		a.semanticError(ErrInvalidOperator)
	}
	if a.tokens[0].kind == symbol && a.lastString() != "EQU" {
		sym := a.tokens[0].val.(string)
		a.symbols[sym] = mix.NewWord(a.self)
		a.patchFixUps(sym)
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
