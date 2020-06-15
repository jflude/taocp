package mix

func (c *Computer) inca(aa Word, i, f, op, m int) int64 {
	switch f {
	case 0: // INC
		if op == INCA || op == INCX {
			c.addAccum(op-INCA, m)
		} else {
			c.addIndex(op-INCA, m)
		}
	case 1: // DEC
		if op == INCA || op == INCX {
			c.addAccum(op-INCA, -m)
		} else {
			c.addIndex(op-INCA, -m)
		}
	case 2: // ENT
		c.Reg[op-INCA] = NewWord(m)
		if m == 0 {
			c.Reg[op-INCA].SetField(FieldSpec(0, 0), aa)
		}
	case 3: // ENN
		c.Reg[op-INCA] = NewWord(-m)
		if m == 0 {
			c.Reg[op-INCA].
				SetField(FieldSpec(0, 0), aa.Negate())
		}
	default:
		panic(ErrInvalidOp)
	}
	return 1
}

func (c *Computer) addIndex(reg, v int) {
	v += c.Reg[reg].Int()
	if v < MinIndex || v > MaxIndex {
		panic(ErrInvalidIndex)
	}
	if v == 0 {
		c.Reg[reg].SetField(FieldSpec(1, 5), 0)
	} else {
		c.Reg[reg] = NewWord(v)
	}
}
