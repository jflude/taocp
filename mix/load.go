package mix

func (c *Computer) lda(aa Word, i, f, op, m int) int {
	v := c.Contents[m].Field(f)
	r := op - LDA
	if r >= I1 && r <= I6 {
		if x := v.Int(); x < MinIndex || x > MaxIndex {
			panic(ErrInvalidIndex)
		}
	}
	c.Reg[r] = v
	return 2
}

func (c *Computer) ldan(aa Word, i, f, op, m int) int {
	v := c.Contents[m].Field(f).Negate()
	r := op - LDAN
	if r >= I1 && r <= I6 {
		if x := v.Int(); x < MinIndex || x > MaxIndex {
			panic(ErrInvalidIndex)
		}
	}
	c.Reg[r] = v
	return 2
}

func (c *Computer) sta(aa Word, i, f, op, m int) int {
	c.Contents[m].SetField(f, c.Reg[op-STA])
	return 2
}