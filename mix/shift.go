package mix

func (c *Computer) sla(aa Word, i, f, op, m int) int64 {
	if m < 0 {
		panic(ErrInvalidOp)
	}
	switch f {
	case 0: // SLA
		c.Reg[A].ShiftLeft(m)
	case 1: // SRA
		c.Reg[A].ShiftRight(m)
	case 2: // SLAX
		c.Reg[A].ShiftLeft(m)
		out := c.Reg[X].ShiftLeft(m)
		c.Reg[A].SetField(FieldSpec(6-m, 5), out)
	case 3: // SRAX
		c.Reg[X].ShiftRight(m)
		out := c.Reg[A].ShiftRight(m)
		c.Reg[X].SetField(FieldSpec(1, m), out)
	case 4: // SLC
		m %= 5
		outA := c.Reg[A].ShiftLeft(m)
		outX := c.Reg[X].ShiftLeft(m)
		c.Reg[A].SetField(FieldSpec(6-m, 5), outX)
		c.Reg[X].SetField(FieldSpec(6-m, 5), outA)
	case 5: // SRC
		m %= 5
		outA := c.Reg[A].ShiftRight(m)
		outX := c.Reg[X].ShiftRight(m)
		c.Reg[A].SetField(FieldSpec(1, m), outX)
		c.Reg[X].SetField(FieldSpec(1, m), outA)
	case 6, 7: // SLB, SRB
		panic(ErrNotImplemented) // TODO: see Section 4.5.2
	default:
		panic(ErrInvalidOp)
	}
	return 2
}
