// Package mixal is an assembler for the MIX computer as simulated by the
// mix package.
package mixal

import (
	"io"

	"github.com/jflude/gnuth/mix"
)

type asmb struct {
	obj         object
	input       string
	tokens      []token
	symbols     map[string]int
	self, count int
	exprOp      int
	exprVal     *mix.Word
	wVal, aa    mix.Word
	i, f, c     int
}

func Assemble(r io.Reader, w io.Writer) error {
	var a asmb
	if err := a.translate(r, parseLine); err != nil {
		return err
	}
	return a.obj.writeCards(w)
}
