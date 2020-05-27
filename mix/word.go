package mix

import (
	"errors"
	"math"
	"strconv"
)

const (
	// The largest and smallest values that are representable by a
	// five byte MIX word.
	MaxWord = 1<<30 - 1
	MinWord = -MaxWord

	// The largest and smallest values that are representable by a
	// two byte MIX index registers and the jump register.
	MaxIndex = 1<<12 - 1
	MinIndex = -MaxIndex

	signBit  = math.MinInt32
	byteSize = 1 << 6
)

var ErrInvalidFieldSpec = errors.New("invalid field specification")

func FieldSpec(left, right int) int {
	return 8*left + right
}

func validateFieldSpec(f int) {
	if f >= len(fields) || fields[f].shift == -1 {
		panic(ErrInvalidFieldSpec)
	}
}

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
	return strconv.Itoa(w.Int())
}

// GoString returns a representation of a MIX word as an unsigned integer.
func (w Word) GoString() string {
	return "0" + strconv.FormatUint(uint64(uint32(w)), 8)
}

// Field returns the value of field f as a MIX word.
func (w Word) Field(f int) Word {
	if f == 5 {
		return w
	}
	validateFieldSpec(f)
	return Word((int32(w) >> fields[f].shift) & fields[f].reg)
}

// SetField changes the field f of the MIX word w to the given value.
func (w *Word) SetField(f int, val Word) {
	if f == 5 {
		*w = val
		return
	}
	validateFieldSpec(f)
	*w = Word((int32(*w) &^ fields[f].mem) |
		((int32(val) << fields[f].shift) & fields[f].mem) |
		(int32(val) & fields[f].sign))
}

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

// Instruction extracts an instruction's address, index, field and opcode
// from a MIX word.
func (w Word) Instruction() (aa Word, i, f, c int) {
	aa = Word((int32(w) >> fields[2].shift) & fields[2].reg)
	i = int((int32(w) & fields[27].mem) >> fields[27].shift)
	f = int((int32(w) & fields[36].mem) >> fields[36].shift)
	c = int((int32(w) & fields[45].mem) >> fields[45].shift)
	return
}

var fields = [...]struct {
	mem   int32 // memory mask
	reg   int32 // register mask
	sign  int32 // sign affected
	shift int   // how much to shift to align receiver and source
}{
	{signBit, signBit, signBit, 0},                              // [0:0]
	{07700000000 | signBit, 00000000077 | signBit, signBit, 24}, // [0:1]
	{07777000000 | signBit, 00000007777 | signBit, signBit, 18}, // [0:2]
	{07777770000 | signBit, 00000777777 | signBit, signBit, 12}, // [0:3]
	{07777777700 | signBit, 00077777777 | signBit, signBit, 6},  // [0:4]
	{07777777777 | signBit, 07777777777 | signBit, signBit, 0},  // [0:5]
	{0, 0, 0, -1},
	{0, 0, 0, -1},
	{0, 0, 0, -1},
	{07700000000, 000000000077, 0, 24}, // [1:1]
	{07777000000, 000000007777, 0, 18}, // [1:2]
	{07777770000, 000000777777, 0, 12}, // [1:3]
	{07777777700, 000077777777, 0, 6},  // [1:4]
	{07777777777, 007777777777, 0, 0},  // [1:5]
	{0, 0, 0, -1},
	{0, 0, 0, -1},
	{0, 0, 0, -1},
	{0, 0, 0, -1},
	{00077000000, 000000000077, 0, 18}, // [2:2]
	{00077770000, 000000007777, 0, 12}, // [2:3]
	{00077777700, 000000777777, 0, 6},  // [2:4]
	{00077777777, 000077777777, 0, 0},  // [2:5]
	{0, 0, 0, -1},
	{0, 0, 0, -1},
	{0, 0, 0, -1},
	{0, 0, 0, -1},
	{0, 0, 0, -1},
	{00000770000, 000000000077, 0, 12}, // [3:3]
	{00000777700, 000000007777, 0, 6},  // [3:4]
	{00000777777, 000000777777, 0, 0},  // [3:5]
	{0, 0, 0, -1},
	{0, 0, 0, -1},
	{0, 0, 0, -1},
	{0, 0, 0, -1},
	{0, 0, 0, -1},
	{0, 0, 0, -1},
	{00000007700, 000000000077, 0, 6}, // [4:4]
	{00000007777, 000000007777, 0, 0}, // [4:5]
	{0, 0, 0, -1},
	{0, 0, 0, -1},
	{0, 0, 0, -1},
	{0, 0, 0, -1},
	{0, 0, 0, -1},
	{0, 0, 0, -1},
	{0, 0, 0, -1},
	{00000000077, 000000000077, 0, 0}, // [5:5]
}
