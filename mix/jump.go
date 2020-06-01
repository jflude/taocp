package mix

func (c *Computer) jump(address int, cond bool) {
	if !cond {
		return
	}
	if address < 0 || address >= MemorySize {
		panic(ErrInvalidAddress)
	}
	c.Reg[J] = NewWord(c.next + 1)
	c.next = address - 1
}

func (c *Computer) jmp(aa Word, i, f, op, m int) int {
	switch f {
	case 0: // JMP
		c.jump(m, true)
	case 1: // JSJ
		if m < 0 || m >= MemorySize {
			panic(ErrInvalidAddress)
		}
		c.next = m - 1
	case 2: // JOV
		c.jump(m, c.Overflow)
		c.Overflow = false
	case 3: // JNV
		c.jump(m, !c.Overflow)
		c.Overflow = false
	case 4: // JL
		c.jump(m, c.Comparison == Less)
	case 5: // JE
		c.jump(m, c.Comparison == Equal)
	case 6: // JG
		c.jump(m, c.Comparison == Greater)
	case 7: // JGE
		c.jump(m, c.Comparison != Less)
	case 8: // JNE
		c.jump(m, c.Comparison != Equal)
	case 9: // JLE
		c.jump(m, c.Comparison != Greater)
	default:
		panic(ErrInvalidInstruction)
	}
	return 1
}
func (c *Computer) ja(aa Word, i, f, op, m int) int {
	r := c.Reg[op-JA].Int()
	switch f {
	case 0: // N
		c.jump(m, r < 0)
	case 1: // Z
		c.jump(m, r == 0)
	case 2: // P
		c.jump(m, r > 0)
	case 3: // NN
		c.jump(m, r >= 0)
	case 4: // NE
		c.jump(m, r != 0)
	case 5: // NP
		c.jump(m, r <= 0)
	default:
		panic(ErrInvalidInstruction)
	}
	return 1
}
