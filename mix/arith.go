package mix

func (c *Computer) add(aa Word, i, f, op, m int) (int64, error) {
	c.checkInterlock(m, m)
	if f == 6 {
		var ov bool
		c.Reg[A], ov = AddFloatWord(c.Reg[A], c.Contents[mBase+m])
		c.Overflow = c.Overflow || ov
		return 4, nil
	}
	c.addReg(A, c.Contents[mBase+m].Field(f).Int())
	return 2, nil
}

func (c *Computer) sub(aa Word, i, f, op, m int) (int64, error) {
	c.checkInterlock(m, m)
	if f == 6 {
		var ov bool
		c.Reg[A], ov = SubFloatWord(c.Reg[A], c.Contents[mBase+m])
		c.Overflow = c.Overflow || ov
		return 4, nil
	}
	c.addReg(A, -c.Contents[mBase+m].Field(f).Int())
	return 2, nil
}

func (c *Computer) mul(aa Word, i, f, op, m int) (int64, error) {
	c.checkInterlock(m, m)
	if f == 6 {
		var ov bool
		c.Reg[A], ov = MulFloatWord(c.Reg[A], c.Contents[mBase+m])
		c.Overflow = c.Overflow || ov
		return 9, nil
	}
	c.Reg[A], c.Reg[X] = MulWord(c.Reg[A],
		c.Contents[mBase+m].Field(f).Int())
	return 10, nil
}

func (c *Computer) div(aa Word, i, f, op, m int) (int64, error) {
	c.checkInterlock(m, m)
	if f == 6 {
		var ov bool
		c.Reg[A], ov = DivFloatWord(c.Reg[A], c.Contents[mBase+m])
		c.Overflow = c.Overflow || ov
		return 11, nil
	}
	var ov bool
	c.Reg[A], c.Reg[X], ov =
		DivWord(c.Reg[A], c.Reg[X], c.Contents[mBase+m].Int())
	c.Overflow = c.Overflow || ov
	return 12, nil
}

func (c *Computer) addReg(reg, v int) {
	var ov bool
	c.Reg[reg], ov = AddWord(c.Reg[reg], v)
	c.Overflow = c.Overflow || ov
}
