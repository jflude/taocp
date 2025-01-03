package mix

import "math"

const (
	excess       = byteSize / 2
	precision    = 4
	expUnderflow = 1
	expOverflow  = 2
)

var nativeExpRatio = int(math.Log2(byteSize))

// AddFloatWord returns the sum of two floating point MIX words, and
// whether overflow or underflow occurred.  See Algorithm A, Section 4.2.1.
func AddFloatWord(u, v Word) (result Word, overflow bool) {
	panic(ErrNotImplemented)
}

func SubFloatWord(u, v Word) (result Word, overflow bool) {
	return AddFloatWord(u, v.Negate())
}

// MulFloatWord returns the product of two floating point MIX words,
// and whether overflow or underflow occurred.  See Algorithm M, Section 4.2.1.
func MulFloatWord(u, v Word) (result Word, overflow bool) {
	panic(ErrNotImplemented)
}

func DivFloatWord(u, v Word) (result Word, overflow bool) {
	panic(ErrNotImplemented)
}

func CompareFloatWord(u, v Word) int {
	panic(ErrNotImplemented)
}

// see the answer to Ex. 14, Section 4.2.1
func FloatToFixed(w Word) (result Word, overflow bool) {
	e, f := w.UnpackFloat()
	if f == 0 {
		return
	}
F1:
	f *= byteSize
	e--
	if abs(f) < 0100000000 {
		goto F1
	}
	n := excess + 4 - e
	if n < 0 {
		overflow = true
		return
	}
	var x int
	for i := 0; i < n; i++ {
		lsb := abs(f) & (byteSize - 1)
		f /= byteSize
		x = (x >> 6) | (lsb << 24)
	}
	if x > 0400000000 || (x == 0400000000 && (abs(f)&1 == 0)) {
		if f < 0 {
			f--
		} else {
			f++
		}
	}
	result = NewWord(f)
	return
}

func FixedToFloat(w Word) (result Word, overflow bool) {
	result, status := normalizeFloat(excess+5, w.Int())
	return result, status != 0
}

// Convert a MIX floating point value to its Go equivalent.
func FloatToNative(w Word) (result float64, status int) {
	e, f := w.UnpackFloat()
	return math.Ldexp(float64(f)/float64(MaxWord+1),
		(e-excess+1)*nativeExpRatio), 0
}

// Convert a Go floating point value to its MIX equivalent.
func NativeToFloat(n float64) (result Word, status int) {
	frac, exp := math.Frexp(n)
	if math.Abs(frac) < 0.5 {
		frac *= 2
		exp--
	}
	return normalizeFloat(excess+exp/nativeExpRatio,
		int(frac*(MaxWord+1)*math.Exp2(float64(exp%nativeExpRatio))))
}

// Normalize the floating point value composed from the specified exponent
// and fraction, rounding to the precision's digits.  Return it as a MIX
// word, and whether overflow or underflow of the exponent occurred during
// the process.  See Algorithm N, Section 4.2.1.
func normalizeFloat(e, f int) (result Word, status int) {
	if f == 0 { // N1
		return
	}
	if abs(f) > MaxWord {
		goto N4
	}
N2:
	if abs(f) >= 0100000000 {
		goto N5
	}
	f *= byteSize // N3
	e--
	goto N2
N4:
	f /= byteSize
	e++
	if abs(f) > MaxWord {
		goto N4
	}
N5:
	tail := abs(f) & (byteSize - 1)
	if tail > excess || (tail == excess && abs(f)&byteSize == 0) {
		if f < 0 {
			f -= byteSize
		} else {
			f += byteSize
		}
		if abs(f) > MaxWord {
			goto N4
		}
	}
	if e < 0 { // N6
		status = expUnderflow
	} else if e >= byteSize {
		status = expOverflow
	}
	result.PackFloat(e, f/byteSize) // N7
	return
}
