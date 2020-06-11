// Package mixal is an assembler for the MIX computer as simulated by the
// mix package.
package mixal

import "io"

type asmb struct {
	obj         object
	self, count int
	input       string
	tokens      []token
	symbols     map[string]int
}

type lineParser func(*asmb, string, string, string)

func Assemble(r io.Reader, w io.Writer) error {
	var a asmb
	if err := a.readProgram(r, parseLine); err != nil {
		return err
	}
	// TODO
	// if err := a.fixUpRefs(); err != nil {
	// 	return err
	// }
	return a.obj.writeCards(w)
}
