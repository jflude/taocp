package mix

// AddFloatWord returns the sum of two floating point MIX words, and
// whether overflow occurred.  See Section 4.2.1.
func AddFloatWord(u, v Word) (result Word, overflow bool) {
	panic(ErrNotImplemented)
}

func SubFloatWord(u, v Word) (result Word, overflow bool) {
	return AddFloatWord(u, v.Negate())
}

func MulFloatWord(u, v Word) (result Word, overflow bool) {
	panic(ErrNotImplemented)
}

func DivFloatWord(u, v Word) (result Word, overflow bool) {
	panic(ErrNotImplemented)
}

func CompareFloatWord(u, v Word) int {
	panic(ErrNotImplemented)
}

func FloatToFixed(w Word) (result Word, overflow bool) {
	panic(ErrNotImplemented)
}

func FixedToFloat(w Word) (result Word, overflow bool) {
	panic(ErrNotImplemented)
}

func normalizeFloat(w Word) (result Word, overflow bool) {
	panic(ErrNotImplemented) // see Algorithm N, Section 4.2.1.
}
