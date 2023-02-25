// Copyright (c) 2020 Justin Flude. All rights reserved.
// Use of this source code is governed by the COPYING.md file.

// Package mixal is an assembler for the MIX 1009 computer as described
// in Donald Knuth's "The Art of Computer Programming" (third edition).
package mixal

import (
	"io"

	"github.com/jflude/taocp/mix"
)

type literal struct {
	sym string
	val mix.Word
}

type asmb struct {
	obj      object
	input    string
	tokens   []token
	symbols  map[string]mix.Word
	fixups   map[string][]int
	literals []literal
	count    int
	self     int
	label    int
	exprOp   int
	exprVal  *mix.Word
	wVal     mix.Word
	aa       mix.Word
	i, f, c  int
}

func Assemble(r io.Reader, w io.Writer, interrupts bool) error {
	var a asmb
	if err := a.translate(r, parseLine); err != nil {
		return err
	}
	return a.obj.writeCards(w, interrupts)
}
