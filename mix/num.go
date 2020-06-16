package mix

import (
	"errors"
	"strconv"
)

var ErrHalted = errors.New("mix: halted")

func (c *Computer) num(aa Word, i, f, op, m int) int64 {
	switch f {
	case 0: // NUM
		var a, x int64
		for i := 1; i <= 5; i++ {
			f := FieldSpec(i, i)
			a = 10*a + int64(c.Reg[A].Field(f).Int())%10
			x = 10*x + int64(c.Reg[X].Field(f).Int())%10
		}
		a = a*100000 + x
		if a > MaxWord {
			a %= MaxWord + 1
			c.Overflow = true
		}
		c.Reg[A].SetField(FieldSpec(1, 5), NewWord(int(a)))
		return 10
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
		return 10
	case 2: // HLT
		panic(ErrHalted)
	case 3: // AND (see Section 4.5.4)
		c.Reg[A] = AndWord(c.Reg[A], abs(c.Contents[mBase+m].Int()))
		return 2
	case 4: // OR (see Section 6.4)
		c.Reg[A] = OrWord(c.Reg[A], abs(c.Contents[mBase+m].Int()))
		return 2
	case 5: // XOR (see Ex. 28, Section 2.5)
		c.Reg[A] = XorWord(c.Reg[A], abs(c.Contents[mBase+m].Int()))
		return 2
	case 6: // FLOT
		panic(ErrNotImplemented) // TODO: see Section 4.2.1
	case 7: // FIX
		panic(ErrNotImplemented) // TODO: see Section 4.2.1
	case 9: // INT
		c.ctrl = !c.ctrl // see Ex. 18, Section 1.4.4
		return 2
	default:
		panic(ErrInvalidOp)
	}
}
