// Copyright (c) 2020 Justin Flude. All rights reserved.
// Use of this source code is governed by the COPYING.md file.
package mix

func (c *Computer) move(aa Word, i, f, op, m int) (int64, error) {
	if f == 0 {
		return 1, nil
	}
	to := c.Reg[I1].Int()
	c.checkAddress(m)
	c.checkAddress(m + f)
	c.checkAddress(to)
	c.checkAddress(to + f)
	c.checkInterlock(to, m+f)
	for j := 0; j < f; j++ {
		c.Contents[mBase+to+j] = c.Contents[mBase+m+j]
	}
	c.Reg[I1] = NewWord(to + f)
	return 1 + 2*int64(f), nil
}
