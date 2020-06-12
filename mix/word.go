package mix

import (
	"fmt"
	"math"
)

const (
	// The largest and smallest values that are representable
	// by a MIX word.
	MaxWord = 1<<30 - 1
	MinWord = -MaxWord

	// The largest and smallest values that are representable
	// by the MIX index registers and the jump register.
	MaxIndex = 1<<12 - 1
	MinIndex = -MaxIndex

	signBit = math.MinInt32
)

// Word represents a MIX machine word, which consists of five 6-bit bytes
// and a sign bit.
type Word int32

// NewWord returns a MIX word with the given integer value.  It panics if
// the value is out of range.
func NewWord(val int) Word {
	if val < MinWord || val > MaxWord {
		panic("out of range")
	}
	if val < 0 {
		return Word(signBit | -val)
	}
	return Word(val)
}

// Int returns the integer value of the given MIX word.
func (w Word) Int() int {
	if int32(w) < 0 {
		return -int(w & MaxWord)

	}
	return int(w & MaxWord)
}

// Sign returns the sign of a MIX word, +1 for positive, -1 for negative.
func (w Word) Sign() int {
	if w&signBit == 0 {
		return 1
	} else {
		return -1
	}
}

// Negate returns the value of a MIX word with its sign inverted.
func (w Word) Negate() Word {
	return w ^ signBit
}

// String returns a representation of the value of a MIX word as a decimal
// number.  Note that a MIX word can have a value of negative zero.
func (w Word) String() string {
	if int32(w) == signBit {
		return "-0"
	}
	return fmt.Sprint(w.Int())
}

// GoString returns a representation of a MIX word as an unsigned integer.
func (w Word) GoString() string {
	s := "+"
	if w&signBit != 0 {
		s = "-"
	}
	return fmt.Sprintf("%s%#011o", s, uint64(uint32(w&^signBit)))
}

// AddWord adds an integer to a MIX word, returning the result as a MIX word,
// and whether overflow occured.  See Section 1.3.1.
func AddWord(w Word, v int) (result Word, overflow bool) {
	v += w.Int()
	if v < MinWord || v > MaxWord {
		overflow = true
		v &= MaxWord
	}
	if v == 0 {
		w.SetField(FieldSpec(1, 5), 0)
	} else {
		w = NewWord(v)
	}
	return w, overflow
}

// SubWord subtracts an integer from a MIX word, returning the result as a
// MIX word, and whether overflow occured.  See Section 1.3.1.
func SubWord(w Word, v int) (result Word, overflow bool) {
	return AddWord(w, -v)
}

// MulWord multiples an integer by a MIX word, returning the product as a
// double-precision MIX word.  See Section 1.3.1.
func MulWord(w Word, v int) (high, low Word) {
	p := int64(w.Int()) * int64(v)
	n := uint64(abs64(p))
	high = NewWord(int((n >> 30) & MaxWord))
	low = NewWord(int(n & MaxWord))
	if p < 0 {
		high = high.Negate()
		low = low.Negate()
	}
	return
}

func abs64(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}

// DivWord divides a double-precision MIX word by an integer, returning the
// quotient and remainder as MIX words, and whether overflow or division by
// zero occured.  See Section 1.3.1.
func DivWord(high, low Word, v int) (quot, rem Word, overflow bool) {
	if v == 0 || abs(high.Int()) >= abs(v) {
		overflow = true
		return
	}
	d := int64(abs(high.Int()))<<30 | int64(abs(low.Int()))
	s := high.Sign()
	if s == -1 {
		d = -d
	}
	q, r := d/int64(v), d%int64(v)
	if (s == -1 && r >= 0) || (s == 1 && r < 0) {
		r = -r
	}
	quot, rem = NewWord(int(q)), NewWord(int(abs64(r)))
	if s == -1 {
		rem = rem.Negate()
	}
	return
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// ShiftLeft shifts the MIX word w left by the given number of bytes,
// returning the bytes shifted out.
func (w *Word) ShiftLeft(count int) Word {
	if count == 0 {
		return *w
	}
	if count >= 5 {
		out := w.Field(FieldSpec(1, 5))
		w.SetField(FieldSpec(1, 5), NewWord(0))
		return out
	}
	if count < 0 {
		panic("invalid shift count")
	}
	out := w.Field(FieldSpec(1, count))
	in := w.Field(FieldSpec(count+1, 5))
	w.SetField(FieldSpec(1, 5-count), in)
	w.SetField(FieldSpec(6-count, 5), 0)
	return out
}

// ShiftRight shifts the MIX word w right by the given number of bytes,
// returning the bytes shifted out.
func (w *Word) ShiftRight(count int) Word {
	if count == 0 {
		return *w
	}
	if count >= 5 {
		out := w.Field(FieldSpec(1, 5))
		w.SetField(FieldSpec(1, 5), NewWord(0))
		return out
	}
	if count < 0 {
		panic("invalid shift count")
	}
	out := w.Field(FieldSpec(6-count, 5))
	in := w.Field(FieldSpec(1, 5-count))
	w.SetField(FieldSpec(count+1, 5), in)
	w.SetField(FieldSpec(1, count), 0)
	return out
}
