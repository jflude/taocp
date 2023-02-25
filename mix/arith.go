// Copyright (c) 2020 Justin Flude. All rights reserved.
// Use of this source code is governed by the COPYING.md file.
package mix

func (c *Computer) add(aa Word, i, f, op, m int) (int64, error) {
	if f == 6 {
		c.callWithOvCheck2(AddFloatWord, m)
		return 4, nil
	}
	c.checkInterlock(m, m)
	c.addReg(A, c.Contents[mBase+m].Field(f).Int())
	return 2, nil
}

func (c *Computer) sub(aa Word, i, f, op, m int) (int64, error) {
	if f == 6 {
		c.callWithOvCheck2(SubFloatWord, m)
		return 4, nil
	}
	c.checkInterlock(m, m)
	c.addReg(A, -c.Contents[mBase+m].Field(f).Int())
	return 2, nil
}

func (c *Computer) mul(aa Word, i, f, op, m int) (int64, error) {
	if f == 6 {
		c.callWithOvCheck2(MulFloatWord, m)
		return 9, nil
	}
	c.checkInterlock(m, m)
	c.Reg[A], c.Reg[X] = MulWord(c.Reg[A],
		c.Contents[mBase+m].Field(f).Int())
	return 10, nil
}

func (c *Computer) div(aa Word, i, f, op, m int) (int64, error) {
	if f == 6 {
		c.callWithOvCheck2(DivFloatWord, m)
		return 11, nil
	}
	var ov bool
	c.checkInterlock(m, m)
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

func (c *Computer) callWithOvCheck2(f func(Word, Word) (Word, bool), m int) {
	var ov bool
	c.checkInterlock(m, m)
	c.Reg[A], ov = f(c.Reg[A], c.Contents[mBase+m])
	c.Overflow = c.Overflow || ov
}
