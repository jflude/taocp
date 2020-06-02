package mix

func (c *Computer) move(aa Word, i, f, op, m int) int64 {
	if f == 0 {
		return 1
	}
	to := c.Reg[I1].Int()
	if m < 0 || m+f >= MemorySize ||
		to < 0 || to+f >= MemorySize {
		panic(ErrInvalidAddress)
	}
	for n := 0; n < f; n++ {
		c.Contents[to+n] = c.Contents[m+n]
	}
	c.Reg[I1] = NewWord(to + f)
	return 1 + 2*int64(f)
}
