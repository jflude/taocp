package mix

func (c *Computer) add(aa Word, i, f, op, m int) int64 {
	if f == 6 {
		panic(ErrNotImplemented)
	}
	c.addAccum(A, c.Contents[m].Field(f).Int())
	return 2
}

func (c *Computer) sub(aa Word, i, f, op, m int) int64 {
	if f == 6 {
		panic(ErrNotImplemented)
	}
	c.addAccum(A, -c.Contents[m].Field(f).Int())
	return 2
}

func (c *Computer) mul(aa Word, i, f, op, m int) int64 {
	if f == 6 {
		panic(ErrNotImplemented)
	}
	c.Reg[A], c.Reg[X] = MulWord(c.Reg[A], c.Contents[m].Field(f).Int())
	return 10
}

func (c *Computer) div(aa Word, i, f, op, m int) int64 {
	if f == 6 {
		panic(ErrNotImplemented)
	}
	var ov bool
	c.Reg[A], c.Reg[X], ov =
		DivWord(c.Reg[A], c.Reg[X], c.Contents[m].Int())
	c.Overflow = c.Overflow || ov
	return 2
}

func (c *Computer) addAccum(reg, v int) {
	var ov bool
	c.Reg[reg], ov = AddWord(c.Reg[reg], v)
	c.Overflow = c.Overflow || ov
}

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

func SubWord(w Word, v int) (result Word, overflow bool) {
	return AddWord(w, -v)
}

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

func abs64(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}
