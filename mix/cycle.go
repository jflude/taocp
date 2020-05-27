package mix

import "strconv"

func (c *Computer) Cycle() (err error) {
	defer func() {
		if err2, ok := recover().(error); ok {
			err = err2
		}
	}()
	if c.next < 0 || c.next >= MemorySize {
		return ErrInvalidAddress
	}
	aa, i, f, op := c.Contents[c.next].Instruction()
	if i > 6 {
		return ErrInvalidInstruction
	}
	m := aa.Int()
	if i > 0 {
		m += c.Reg[i].Int()
		if m < MinIndex || m > MaxIndex {
			return ErrInvalidIndex
		}
	}
	var t int
	switch {
	case op == NOP:
		t = 1
	case op == ADD:
		if f == 6 {
			return ErrNotImplemented
		}
		c.addAccum(A, c.Contents[m].Field(f).Int())
		t = 2
	case op == SUB:
		if f == 6 {
			return ErrNotImplemented
		}
		c.addAccum(A, -c.Contents[m].Field(f).Int())
		t = 2
	case op == MUL:
		if f == 6 {
			return ErrNotImplemented
		}
		v := int64(c.Reg[A].Int()) * int64(c.Contents[m].Field(f).Int())
		ax := uint64(abs64(v))
		c.Reg[A] = NewWord(int((ax >> 30) & MaxWord))
		c.Reg[X] = NewWord(int(ax & MaxWord))
		if v < 0 {
			c.Reg[A] = c.Reg[A].Negate()
			c.Reg[X] = c.Reg[X].Negate()
		}
		t = 10
	case op == DIV:
		if f == 6 {
			return ErrNotImplemented
		}
		v := c.Contents[m].Int()
		if v == 0 || abs(c.Reg[A].Int()) >= abs(v) {
			c.Overflow = true
		} else {
			d := int64(abs(c.Reg[A].Int()))<<30 |
				int64(abs(c.Reg[X].Int()))
			sign := c.Reg[A].Sign()
			if sign == -1 {
				d = -d
			}
			q, r := d/int64(v), d%int64(v)
			c.Reg[A] = NewWord(int(q))
			if (sign == -1 && r >= 0) || (sign == 1 && r < 0) {
				r = -r
			}
			c.Reg[X] = NewWord(int(abs64(r)))
			if sign == -1 {
				c.Reg[X] = c.Reg[X].Negate()
			}

		}
		t = 12
	case op == NUM:
		switch f {
		case 0: // NUM
			var a, x int
			for i := 1; i <= 5; i++ {
				f := FieldSpec(i, i)
				a = 10*a + (c.Reg[A].Field(f).Int() & 07)
				x = 10*x + (c.Reg[X].Field(f).Int() & 07)
			}
			v := a*100000 + x
			c.Reg[A].SetField(FieldSpec(1, 5), NewWord(v))
		case 1: // CHAR
			v := strconv.Itoa(c.Reg[A].Field(FieldSpec(1, 5)).Int())
			if l := len(v); l < 10 {
				v = "000000000"[:10-l] + v
			}
			for i := 0; i < 5; i++ {
				f := FieldSpec(i+1, i+1)
				a := NewWord(utf2mix[rune(v[i])])
				x := NewWord(utf2mix[rune(v[i+5])])
				c.Reg[A].SetField(f, a)
				c.Reg[X].SetField(f, x)
			}
		case 2: // HLT
			err = ErrHalted
		default:
			return ErrInvalidInstruction
		}
		t = 1
	case op == SLA:
		if m < 0 {
			return ErrInvalidInstruction
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
			m %= 6
			outA := c.Reg[A].ShiftLeft(m)
			outX := c.Reg[X].ShiftLeft(m)
			c.Reg[A].SetField(FieldSpec(6-m, 5), outX)
			c.Reg[X].SetField(FieldSpec(6-m, 5), outA)
		case 5: // SRC
			m %= 6
			outA := c.Reg[A].ShiftRight(m)
			outX := c.Reg[X].ShiftRight(m)
			c.Reg[A].SetField(FieldSpec(1, m), outX)
			c.Reg[X].SetField(FieldSpec(1, m), outA)
		default:
			return ErrInvalidInstruction
		}
		t = 2
	case op == MOVE:
		i := c.Reg[I1].Int()
		if m < 0 || m+f >= MemorySize || i < 0 || i+f >= MemorySize {
			return ErrInvalidAddress
		}
		for n := 0; n < f; n++ {
			c.Contents[i+n] = c.Contents[m+n]
		}
		c.Reg[I1] = NewWord(i + f)
		t = 1 + 2*int(f)
	case op >= LDA && op <= LDX:
		v := c.Contents[m].Field(f)
		r := op - LDA
		if r >= I1 && r <= I6 {
			if x := v.Int(); x < MinIndex || x > MaxIndex {
				return ErrInvalidIndex
			}
		}
		c.Reg[r] = v
		t = 2
	case op >= LDAN && op <= LDXN:
		v := c.Contents[m].Field(f).Negate()
		r := op - LDAN
		if r >= I1 && r <= I6 {
			if x := v.Int(); x < MinIndex || x > MaxIndex {
				return ErrInvalidIndex
			}
		}
		c.Reg[r] = v
		t = 2
	case op >= STA && op <= STZ:
		c.Contents[m].SetField(f, c.Reg[op-STA])
		t = 2
	case op == JBUS:
		if f >= len(c.Devices) {
			return ErrInvalidInstruction
		}
		c.jump(m, c.Devices[f].BusyUntil(c.elapsed) != 0)
		t = 1
	case op == IOC:
		if f >= len(c.Devices) {
			return ErrInvalidInstruction
		}
		c.waitBusy(f)
		if err := c.Devices[f].Control(m); err != nil {
			return err
		}
		t = 1
	case op == IN:
		if f >= len(c.Devices) {
			return ErrInvalidInstruction
		}
		n := m + c.Devices[f].BlockSize()
		if m < 0 || n >= MemorySize {
			return ErrInvalidAddress
		}
		c.waitBusy(f)
		if err := c.Devices[f].Read(c.Contents[m:n]); err != nil {
			return err
		}
		t = 1
	case op == OUT:
		if f >= len(c.Devices) {
			return ErrInvalidInstruction
		}
		n := m + c.Devices[f].BlockSize()
		if m < 0 || n >= MemorySize {
			return ErrInvalidAddress
		}
		c.waitBusy(f)
		if err := c.Devices[f].Write(c.Contents[m:n]); err != nil {
			return err
		}
		t = 1
	case op == JRED:
		if f >= len(c.Devices) {
			return ErrInvalidInstruction
		}
		c.jump(m, c.Devices[f].BusyUntil(c.elapsed) == 0)
		t = 1
	case op == JMP:
		switch f {
		case 0: // JMP
			c.jump(m, true)
		case 1: // JSJ
			c.next = m
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
			return ErrInvalidInstruction
		}
		t = 1
	case op >= JA && op <= JX:
		reg := c.Reg[op-JA].Int()
		switch f {
		case 0: // N
			c.jump(m, reg < 0)
		case 1: // Z
			c.jump(m, reg == 0)
		case 2: // P
			c.jump(m, reg > 0)
		case 3: // NN
			c.jump(m, reg >= 0)
		case 4: // NE
			c.jump(m, reg != 0)
		case 5: // NP
			c.jump(m, reg <= 0)
		default:
			return ErrInvalidInstruction
		}
		t = 1
	case op == INCA || op == INCX:
		switch f {
		case 0: // INC
			c.addAccum(op-INCA, m)
		case 1: // DEC
			c.addAccum(op-INCA, -m)
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
			return ErrInvalidInstruction
		}
		t = 1
	case op >= INC1 && op <= INC6:
		switch f {
		case 0: // INC
			if err := c.addIndex(op-INCA, m); err != nil {
				return err
			}
		case 1: // DEC
			if err := c.addIndex(op-INCA, -m); err != nil {
				return err
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
			return ErrInvalidInstruction
		}
		t = 1
	case op >= CMPA && op <= CMPX:
		reg := c.Reg[op-CMPA].Field(f).Int()
		mem := c.Contents[m].Field(f).Int()
		if reg < mem {
			c.Comparison = Less
		} else if reg > mem {
			c.Comparison = Greater
		} else {
			c.Comparison = Equal
		}
		t = 2
	}
	c.elapsed += int64(t)
	c.next++
	return err
}

func (c *Computer) addAccum(reg, v int) {
	v += c.Reg[reg].Int()
	if v < MinWord || v > MaxWord {
		c.Overflow = true
		v &= MaxWord
	}
	if v == 0 {
		c.Reg[reg].SetField(FieldSpec(1, 5), 0)
	} else {
		c.Reg[reg] = NewWord(v)
	}
}

func (c *Computer) addIndex(reg, v int) error {
	v += c.Reg[reg].Int()
	if v < MinIndex || v > MaxIndex {
		return ErrInvalidIndex
	}
	if v == 0 {
		c.Reg[reg].SetField(FieldSpec(1, 5), 0)
	} else {
		c.Reg[reg] = NewWord(v)
	}
	return nil
}

func (c *Computer) jump(address int, cond bool) {
	if cond {
		c.Reg[J] = NewWord(c.next + 1)
		c.next = address
	}
}

func (c *Computer) waitBusy(unit int) {
	if until := c.Devices[unit].BusyUntil(c.elapsed); until > 0 {
		c.elapsed = until
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func abs64(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}
