// Copyright (c) 2020 Justin Flude. All rights reserved.
// Use of this source code is governed by the COPYING.md file.
package mix

import "testing"

func TestWord(t *testing.T) {
	w := NewWord(NewWord(-3).Int() + 4)
	checkWord(t, w, NewWord(1))

	w = w.Negate()
	x := w.Field(Spec(5, 5))
	checkWord(t, x, NewWord(1))

	x = w.Field(Spec(0, 1))
	checkWord(t, x, NewWord(0).Negate())

	y := NewWord(w.Int() - 05555555555)
	checkWord(t, y, NewWord(-05555555556))

	x = y.Field(Spec(0, 2))
	checkWord(t, x, NewWord(-05555))

	x = y.Field(Spec(1, 2))
	checkWord(t, x, NewWord(05555))

	w = NewWord(-1)
	w.SetField(Spec(1, 2), NewWord(0737))
	checkWord(t, w, NewWord(-0737000001))

	w.SetField(Spec(3, 4), NewWord(04567))
	checkWord(t, w, NewWord(-0737456701))

	w = NewWord(-1)
	w.SetField(Spec(0, 2), NewWord(0777))
	checkWord(t, w, NewWord(0777000001))

	w = NewWord(1)
	w.SetField(Spec(0, 4), NewWord(-07777))
	checkWord(t, w, NewWord(-0777701))

	w = NewWord(1)
	w.SetField(Spec(0, 0), NewWord(-07777777777))
	checkWord(t, w, NewWord(-1))

	var hi Word
	lo := NewWord(-01001234567)
	ShiftBitsLeft(&hi, &lo, 12)
	checkWord(t, hi, NewWord(01001))
	checkWord(t, lo, NewWord(-02345670000))

	hi = 0
	lo = NewWord(-01001234567)
	ShiftBitsLeft(&hi, &lo, 33)
	checkWord(t, hi, NewWord(012345670))
	checkWord(t, lo, NewWord(0).Negate())

	hi = NewWord(-01001234567)
	lo = 0
	ShiftBitsRight(&hi, &lo, 6)
	checkWord(t, hi, NewWord(-010012345))
	checkWord(t, lo, NewWord(06700000000))

	hi = NewWord(-01001234567)
	lo = 0
	ShiftBitsRight(&hi, &lo, 36)
	checkWord(t, hi, NewWord(0).Negate())
	checkWord(t, lo, NewWord(010012345))

	hi = NewWord(-0123456701)
	lo = NewWord(-02345670123)
	ShiftBitsRight(&hi, &lo, 60)
	checkWord(t, hi, NewWord(0).Negate())
	checkWord(t, lo, NewWord(0).Negate())

	hi = NewWord(-0123456701)
	lo = NewWord(-02345670123)
	RotateBitsLeft(&hi, &lo, 12)
	checkWord(t, hi, NewWord(-04567012345))
	checkWord(t, lo, NewWord(-06701230123))

	hi = NewWord(-0123456701)
	lo = NewWord(-02345670123)
	RotateBitsRight(&hi, &lo, 6)
	checkWord(t, hi, NewWord(-02301234567))
	checkWord(t, lo, NewWord(-0123456701))

	hi = NewWord(-0123456701)
	lo = NewWord(-02345670123)
	RotateBitsRight(&hi, &lo, 120)
	checkWord(t, hi, NewWord(-0123456701))
	checkWord(t, lo, NewWord(-02345670123))

	w = NewWord(-01001234567)
	aa, i, f, c := w.UnpackOp()
	if aa.Int() != -01001 || i != 023 || f != 045 || c != 067 {
		t.Errorf("got: %#o, %#o, %#o, %#o, want: -01001, 023, 045, 067",
			aa.Int(), i, f, c)
	}

	w.PackOp(aa, i, f, c)
	if w.Int() != -01001234567 {
		t.Errorf("got: %#o, want: -01001234567", w.Int())
	}

	if s := NewWord(0).Negate().String(); s != "-0" {
		t.Errorf(`got: %q, want: "-0"`, s)
	}

	valid := true
	func() {
		defer func() {
			if err := recover(); err != nil {
				valid = false
			}
		}()
		w.Field(Spec(1, 6))
	}()
	if valid {
		t.Error("invalid Spec did not panic")
	}
}

func checkWord(t *testing.T, got, want Word) {
	t.Helper()
	if int32(got) != int32(want) {
		t.Errorf("got: %#v (%v), want: %#v", got, got, want)
	}
}
