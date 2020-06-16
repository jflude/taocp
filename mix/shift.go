package mix

func (c *Computer) sla(aa Word, i, f, op, m int) int64 {
	if m < 0 {
		panic(ErrInvalidOp)
	}
	switch f {
	case 0: // SLA
		c.Reg[A].ShiftBytesLeft(m)
	case 1: // SRA
		c.Reg[A].ShiftBytesRight(m)
	case 2: // SLAX
		c.Reg[A].ShiftBytesLeft(m)
		out := c.Reg[X].ShiftBytesLeft(m)
		c.Reg[A].SetField(FieldSpec(6-m, 5), out)
	case 3: // SRAX
		c.Reg[X].ShiftBytesRight(m)
		out := c.Reg[A].ShiftBytesRight(m)
		c.Reg[X].SetField(FieldSpec(1, m), out)
	case 4: // SLC
		m %= 5
		outA := c.Reg[A].ShiftBytesLeft(m)
		outX := c.Reg[X].ShiftBytesLeft(m)
		c.Reg[A].SetField(FieldSpec(6-m, 5), outX)
		c.Reg[X].SetField(FieldSpec(6-m, 5), outA)
	case 5: // SRC
		m %= 5
		outA := c.Reg[A].ShiftBytesRight(m)
		outX := c.Reg[X].ShiftBytesRight(m)
		c.Reg[A].SetField(FieldSpec(1, m), outX)
		c.Reg[X].SetField(FieldSpec(1, m), outA)
	case 6: // SLB (see Section 4.5.2)
		c.Reg[A].ShiftBitsLeft(m)
		out := c.Reg[X].ShiftBitsLeft(m)
		c.Reg[A] = OrWord(c.Reg[A], out.Int())
	case 7: // SRB (see Section 4.5.2)
		c.Reg[X].ShiftBitsRight(m)
		out := c.Reg[A].ShiftBitsRight(m)
		c.Reg[X] = OrWord(c.Reg[X], out.Int()<<(30-m))
	default:
		panic(ErrInvalidOp)
	}
	return 2
}
