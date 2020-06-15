package mix

func (c *Computer) add(aa Word, i, f, op, m int) int64 {
	if f == 6 {
		panic(ErrNotImplemented)
	}
	c.addAccum(A, c.Contents[mBase+m].Field(f).Int())
	return 2
}

func (c *Computer) sub(aa Word, i, f, op, m int) int64 {
	if f == 6 {
		panic(ErrNotImplemented)
	}
	c.addAccum(A, -c.Contents[mBase+m].Field(f).Int())
	return 2
}

func (c *Computer) mul(aa Word, i, f, op, m int) int64 {
	if f == 6 {
		panic(ErrNotImplemented)
	}
	c.Reg[A], c.Reg[X] = MulWord(c.Reg[A],
		c.Contents[mBase+m].Field(f).Int())
	return 10
}

func (c *Computer) div(aa Word, i, f, op, m int) int64 {
	if f == 6 {
		panic(ErrNotImplemented)
	}
	var ov bool
	c.Reg[A], c.Reg[X], ov =
		DivWord(c.Reg[A], c.Reg[X], c.Contents[mBase+m].Int())
	c.Overflow = c.Overflow || ov
	return 12
}

func (c *Computer) addAccum(reg, v int) {
	var ov bool
	c.Reg[reg], ov = AddWord(c.Reg[reg], v)
	c.Overflow = c.Overflow || ov
}
