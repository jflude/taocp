package mix

import "fmt"

func (c *Computer) printTrace(m, next int) {
	var ov, ci string
	if c.Overflow {
		ov = "Y"
	} else {
		ov = "N"
	}
	if c.Comparison == Less {
		ci = "<"
	} else if c.Comparison == Greater {
		ci = ">"
	} else {
		ci = "="
	}
	fmt.Printf(" A: %10v (%#v)  OP: %4d: %s\n X: %10v (%#v)  OV: %s\n"+
		"I1:       %4v (%#v)  CI: %s\n",
		c.Reg[A], c.Reg[A], next,
		Disassemble(c.Contents[mBase+next]),
		c.Reg[X], c.Reg[X], ov, c.Reg[I1], c.Reg[I1], ci)
	for i := 2; i <= 6; i, m = i+1, m+1 {
		if c.validAddress(m) {
			fmt.Printf("I%d:       %4v (%#v)      %4d: %#v\n",
				i, c.Reg[i], c.Reg[i], m,
				c.Contents[mBase+m])
		} else {
			fmt.Printf("I%d:       %4v (%#v)      %4d: ?\n",
				i, c.Reg[i], c.Reg[i], m)
		}
	}
	if m >= 0 && m < MemorySize {
		fmt.Printf("J:        %4v (%#v)      %4d: %#v\n",
			c.Reg[J], c.Reg[J], m, c.Contents[mBase+m])
	} else {
		fmt.Printf("J:        %4v (%#v)      %4d: ?\n",
			c.Reg[J], c.Reg[J], m)
	}
	b := make([]byte, len(c.Devices))
	for i := 0; i < 20; i++ {
		if c.isBusy(i) {
			b[i] += c.Devices[i].Name()[0]
		} else {
			b[i] += '.'
		}
	}
	fmt.Printf("Device:  %s   Elapsed: %d\n\n", string(b), c.elapsed)
}
