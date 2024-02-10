package mixal

import "github.com/jflude/taocp/mix"

func (a *asmb) parseEQU() {
	if a.wVal = 0; !a.parseWValue() {
		a.syntaxError()
	}
	sym := a.tokens[0].val.(string)
	a.symbols[sym] = a.wVal
	a.patchFixUps(sym)
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
	a.input = ""
	a.emit(w[0])
}

func (a *asmb) parseEND() {
	if a.wVal = 0; !a.parseWValue() {
		a.syntaxError()
	}
	a.obj.start = a.wVal.Int()
	for _, lit := range a.literals {
		a.symbols[lit.sym] = mix.NewWord(a.self)
		a.patchFixUps(lit.sym)
		a.emit(lit.val)
	}
	for k := range a.fixups {
		a.symbols[k] = mix.NewWord(a.self)
		a.patchFixUps(k)
		a.emit(0)
	}
}
