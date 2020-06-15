package mix

import (
	"errors"
	"strconv"
)

var ErrHalted = errors.New("halted")

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
		panic(ErrHalted)
	case 7: // INT
		panic(ErrNotImplemented) // TODO: see Ex. 18, Section 1.4.4
	default:
		panic(ErrInvalidInstruction)
	}
	return 1
}
