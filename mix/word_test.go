package mix

import "testing"

func TestWord(t *testing.T) {
	w := NewWord(NewWord(-3).Int() + 4)
	checkWord(t, w, NewWord(1))

	w = w.Negate()
	x := w.Field(FieldSpec(5, 5))
	checkWord(t, x, NewWord(1))

	x = w.Field(FieldSpec(0, 1))
	checkWord(t, x, NewWord(0).Negate())

	y := NewWord(w.Int() - 05555555555)
	checkWord(t, y, NewWord(-05555555556))

	x = y.Field(FieldSpec(0, 2))
	checkWord(t, x, NewWord(-05555))

	x = y.Field(FieldSpec(1, 2))
	checkWord(t, x, NewWord(05555))

	w = NewWord(-1)
	w.SetField(FieldSpec(1, 2), NewWord(0737))
	checkWord(t, w, NewWord(-0737000001))

	w.SetField(FieldSpec(3, 4), NewWord(04567))
	checkWord(t, w, NewWord(-0737456701))

	w = NewWord(-1)
	w.SetField(FieldSpec(0, 2), NewWord(0777))
	checkWord(t, w, NewWord(0777000001))

	w = NewWord(1)
	w.SetField(FieldSpec(0, 4), NewWord(-07777))
	checkWord(t, w, NewWord(-0777701))

	w = NewWord(1)
	w.SetField(FieldSpec(0, 0), NewWord(-07777777777))
	checkWord(t, w, NewWord(-1))

	w = NewWord(-01001234567)
	out := w.ShiftLeft(2)
	checkWord(t, w, NewWord(-02345670000))
	if out != 01001 {
		t.Errorf("got: %#o, want: 01001", out)
	}

	w = NewWord(-01001234567)
	out = w.ShiftRight(2)
	checkWord(t, w, NewWord(-0100123))
	if out != 04567 {
		t.Errorf("got: %#o, want: 04567", out)
	}

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
		w.Field(FieldSpec(1, 6))
	}()
	if valid {
		t.Error("invalid FieldSpec did not panic")
	}
}

func checkWord(t *testing.T, got, want Word) {
	t.Helper()
	if int32(got) != int32(want) {
		t.Errorf("got: %#v (%v), want: %#v", got, got, want)
	}
}
