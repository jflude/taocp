package mix

import "strconv"

func (c *Computer) num(aa Word, i, f, op, m int) int {
	switch f {
	case 0: // NUM
		var a, x int
		for i := 1; i <= 5; i++ {
			f := FieldSpec(i, i)
			a = 10*a + (c.Reg[A].Field(f).Int() % 10)
			x = 10*x + (c.Reg[X].Field(f).Int() % 10)
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
		panic(ErrHalted)
	//case 3: // TODO: INT (f == 7), see Ex. 18, Section 1.4.4
	default:
		panic(ErrInvalidInstruction)
	}
	return 1
}
