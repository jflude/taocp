package mix

func (c *Computer) cmpa(aa Word, i, f, op, m int) (int64, error) {
	c.checkInterlock(m, m)
	if op == CMPA && f == 6 {
		c.Comparison = CompareFloatWord(c.Reg[A], c.Contents[mBase+m])
		return 4, nil
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
	return 2, nil
}
