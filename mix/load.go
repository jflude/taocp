// Copyright (c) 2020 Justin Flude. All rights reserved.
// Use of this source code is governed by the COPYING.md file.
package mix

func (c *Computer) lda(aa Word, i, f, op, m int) int64 {
	c.checkInterlock(m, m)
	v := c.Contents[mBase+m].Field(f)
	r := op - LDA
	if r >= I1 && r <= I6 {
		if x := v.Int(); x < MinIndex || x > MaxIndex {
			panic(ErrInvalidIndex)
		}
	}
	c.Reg[r] = v
	return 2
}

func (c *Computer) ldan(aa Word, i, f, op, m int) int64 {
	c.checkInterlock(m, m)
	v := c.Contents[mBase+m].Field(f).Negate()
	r := op - LDAN
	if r >= I1 && r <= I6 {
		if x := v.Int(); x < MinIndex || x > MaxIndex {
			panic(ErrInvalidIndex)
		}
	}
	c.Reg[r] = v
	return 2
}

func (c *Computer) sta(aa Word, i, f, op, m int) int64 {
	c.checkInterlock(m, m)
	c.Contents[mBase+m].SetField(f, c.Reg[op-STA])
	return 2
}
