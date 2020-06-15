package mix

func (c *Computer) cmpa(aa Word, i, f, op, m int) int64 {
	if op == CMPA && f == 6 {
		panic(ErrNotImplemented)
	}
	reg := c.Reg[op-CMPA].Field(f).Int()
	mem := c.Contents[mBase+m].Field(f).Int()
	if reg < mem {
		c.Comparison = Less
	} else if reg > mem {
		c.Comparison = Greater
	} else {
		c.Comparison = Equal
	}
	return 2
}
