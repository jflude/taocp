package mix

import "errors"

var ErrInvalidSpec = errors.New("mix: invalid field specification")

// Spec returns the integer equivalent of a MIX field specification.
func Spec(left, right int) int {
	return 8*left + right
}

func checkSpec(f int) {
	if f >= len(fields) || fields[f].shift == -1 {
		panic(ErrInvalidSpec)
	}
}

// Field returns the value of MIX word w's field f, as a MIX word.
func (w Word) Field(f int) Word {
	if f == 5 {
		return w
	}
	checkSpec(f)
	return Word((int32(w) >> fields[f].shift) & fields[f].reg)
}

// SetField changes the field f of the MIX word w to the given value.
func (w *Word) SetField(f int, val Word) {
	if f == 5 {
		*w = val
		return
	}
	checkSpec(f)
	*w = Word((int32(*w) &^ fields[f].mem) |
		((int32(val) << fields[f].shift) & fields[f].mem) |
		(int32(val) & fields[f].sign))
}

// PackOp composes a MIX word from a MIX instruction's address, index, field
// and opcode.
func (w *Word) PackOp(aa Word, i, f, c int) {
	*w = Word((int32(aa) & fields[02].sign) |
		(int32(aa) << fields[02].shift & fields[02].mem) |
		(int32(i) << fields[033].shift & fields[033].mem) |
		(int32(f) << fields[044].shift & fields[044].mem) |
		(int32(c) << fields[055].shift & fields[055].mem))
}

// UnpackOp extracts a MIX instruction's address, index, field and opcode
// from a MIX word.
func (w Word) UnpackOp() (aa Word, i, f, c int) {
	aa = Word((int32(w) >> fields[02].shift) & fields[02].reg)
	i = int((int32(w) & fields[033].mem) >> fields[033].shift)
	f = int((int32(w) & fields[044].mem) >> fields[044].shift)
	c = int((int32(w) & fields[055].mem) >> fields[055].shift)
	return
}

// PackFloat composes a MIX word from a floating point number's exponent
// and fraction.
func (w *Word) PackFloat(e, f int) {
	*w = Word(0) // TODO
}

// UnpackFloat extracts the exponent and fraction of a floating point number
// in a MIX word.
func (w Word) UnpackFloat() (e, f int) {
	e = int((int32(w) & fields[011].mem) >> fields[011].shift) // ??
	f = int(0) // TODO
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
