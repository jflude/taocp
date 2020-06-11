package mixal

import "github.com/jflude/gnuth/mix"

func (a *asmb) parseEQU() {
	if !a.matchWValue() {
		a.syntaxError()
	}
	// TODO
}

func (a *asmb) parseORIG() {
	if !a.matchWValue() {
		a.syntaxError()
	}
	// TODO
}

func (a *asmb) parseCON() {
	if !a.matchWValue() {
		a.syntaxError()
	}
	// TODO
}

func (a *asmb) parseALF() {
	w, err := mix.ConvertToMIX(a.input)
	if err != nil {
		parseError(err, a.input)
	}
	a.emit(w[0])
}

func (a *asmb) parseEND() {
	if !a.matchWValue() {
		a.syntaxError()
	}
	// TODO
}
