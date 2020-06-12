package mixal

import "github.com/jflude/gnuth/mix"

func (a *asmb) parseEQU() {
	if a.wVal = 0; !a.parseWValue() {
		a.syntaxError()
	}
	a.symbols[a.tokens[0].val.(string)] = a.wVal.Int()
}

func (a *asmb) parseORIG() {
	if a.wVal = 0; !a.parseWValue() {
		a.syntaxError()
	}
	a.self = a.wVal.Int()
	a.newSegment(a.self)
}

func (a *asmb) parseCON() {
	if a.wVal = 0; !a.parseWValue() {
		a.syntaxError()
	}
	a.emit(a.wVal)
}

func (a *asmb) parseALF() {
	w, err := mix.ConvertToMIX(a.input)
	if err != nil {
		parseError(err, a.input)
	}
	a.emit(w[0])
}

func (a *asmb) parseEND() {
	if a.wVal = 0; !a.parseWValue() {
		a.syntaxError()
	}
	a.obj.start = a.wVal.Int()
	// TODO: fix-up and emit literals and undefined symbols
}
