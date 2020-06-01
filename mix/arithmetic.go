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
	v := int64(c.Reg[A].Int()) *
		int64(c.Contents[m].Field(f).Int())
	ax := uint64(abs64(v))
	c.Reg[A] = NewWord(int((ax >> 30) & MaxWord))
	c.Reg[X] = NewWord(int(ax & MaxWord))
	if v < 0 {
		c.Reg[A] = c.Reg[A].Negate()
		c.Reg[X] = c.Reg[X].Negate()
	}
	return 10
}

func (c *Computer) div(aa Word, i, f, op, m int) int64 {
	if f == 6 {
		panic(ErrNotImplemented)
	}
	v := c.Contents[m].Int()
	if v == 0 || abs(c.Reg[A].Int()) >= abs(v) {
		c.Overflow = true
	} else {
		d := int64(abs(c.Reg[A].Int()))<<30 |
			int64(abs(c.Reg[X].Int()))
		sign := c.Reg[A].Sign()
		if sign == -1 {
			d = -d
		}
		q, r := d/int64(v), d%int64(v)
		c.Reg[A] = NewWord(int(q))
		if (sign == -1 && r >= 0) || (sign == 1 && r < 0) {
			r = -r
		}
		c.Reg[X] = NewWord(int(abs64(r)))
		if sign == -1 {
			c.Reg[X] = c.Reg[X].Negate()
		}
	}
	return 2
}

func (c *Computer) addAccum(reg, v int) {
	v += c.Reg[reg].Int()
	if v < MinWord || v > MaxWord {
		c.Overflow = true
		v &= MaxWord
	}
	if v == 0 {
		c.Reg[reg].SetField(FieldSpec(1, 5), 0)
	} else {
		c.Reg[reg] = NewWord(v)
	}
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
