package mix

import "fmt"

func (c *Computer) printTrace(m, next int) {
	var ov, ci, ctrl string
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
	if c.ctrl {
		ctrl = fmt.Sprintf("CTRL: %d", c.pending.Len())
	}
	asm := Disassemble(c.Contents[mBase+next])
	fmt.Fprintf(c.Tracer,
		"\014_______________________________________________________\n"+
			" A: %10v (%#v)   OP: %4d%s %s\n"+
			" X: %10v (%#v)   OV: %s CI: %s %s\n",
		c.Reg[A], c.Reg[A], next, c.lockChar(next), asm,
		c.Reg[X], c.Reg[X], ov, ci, ctrl)
	_, _, _, op := c.Contents[mBase+c.next].UnpackOp()
	for i := 1; i <= 6; i, m = i+1, m+1 {
		mLabel := "  "
		if i == 4 {
			mLabel = "M:"
		}
		lock, mem := c.memoryStyle(m, op)
		fmt.Fprintf(c.Tracer, "I%d:       %4v (%#v)    %s%5d%s %s\n",
			i, c.Reg[i], c.Reg[i], mLabel, m, lock, mem)
	}
	lock, mem := c.memoryStyle(m, op)
	fmt.Fprintf(c.Tracer, " J:       %4v (%#v)      %5d%s %s\n",
		c.Reg[J], c.Reg[J], m, lock, mem)
	m++
	b := make([]byte, len(c.Devices))
	var devMask uint
	for i := 0; i < DeviceCount; i++ {
		if c.isBusy(i) {
			b[i] += c.Devices[i].Name()[0]
			devMask |= 1
		} else {
			b[i] += '.'
		}
		devMask <<= 1
	}
	flipped := ":"
	if devMask != c.lastDevMask {
		flipped = "!"
		c.lastDevMask = devMask
	}
	idled := ":"
	if c.Idle != c.lastIdle {
		idled = "!"
		c.lastIdle = c.Idle
	}
	lock, mem = c.memoryStyle(m, op)
	fmt.Fprintf(c.Tracer, "  Unit%s %s      %5d%s %s\n"+
		"  Idle%s          %11du    Elapsed: %11du\n",
		flipped, string(b), m, lock, mem, idled, c.Idle, c.Elapsed)
}

func (c *Computer) memoryStyle(m, op int) (lock, mem string) {
	if c.validAddress(m) {
		if op < JRED || op > JX {
			mem = c.Contents[mBase+m].GoString()
		} else {
			mem = Disassemble(c.Contents[mBase+m])
		}
	} else {
		mem = "?"
	}
	lock = c.lockChar(m)
	return
}

func (c *Computer) lockChar(m int) string {
	if c.validAddress(m) && c.attributes[mBase+m].lockUntil > c.Elapsed {
		return "#"
	} else {
		return ":"
	}
}
