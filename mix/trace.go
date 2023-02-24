// Copyright (c) 2020 Justin Flude. All rights reserved.
// Use of this source code is governed by the COPYING.md file.
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
	fmt.Fprintf(c.Tracer, "\014\n A: %10v (%#v)   OP: %4d: %s\n"+
		" X: %10v (%#v)   OV: %s CI: %s %s\n"+
		"I1:       %4v (%#v)                   M\n",
		c.Reg[A], c.Reg[A], next, Disassemble(c.Contents[mBase+next]),
		c.Reg[X], c.Reg[X], ov, ci, ctrl, c.Reg[I1], c.Reg[I1])
	for i := 2; i <= 6; i, m = i+1, m+1 {
		if c.validAddress(m) {
			fmt.Fprintf(c.Tracer,
				"I%d:       %4v (%#v)      %5d: %#v\n",
				i, c.Reg[i], c.Reg[i], m,
				c.Contents[mBase+m])
		} else {
			fmt.Fprintf(c.Tracer,
				"I%d:       %4v (%#v)      %5d: ?\n",
				i, c.Reg[i], c.Reg[i], m)
		}
	}
	if c.validAddress(m) {
		fmt.Fprintf(c.Tracer, "J:        %4v (%#v)      %5d: %#v\n",
			c.Reg[J], c.Reg[J], m, c.Contents[mBase+m])
	} else {
		fmt.Fprintf(c.Tracer, "J:        %4v (%#v)      %5d: ?\n",
			c.Reg[J], c.Reg[J], m)
	}
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
	devChanged := ":"
	if devMask != c.lastDevMask {
		devChanged = "!"
		c.lastDevMask = devMask
	}
	idleChanged := ":"
	if c.Idle != c.lastIdle {
		idleChanged = "!"
		c.lastIdle = c.Idle
	}
	fmt.Fprintf(c.Tracer,
		"Device%s %s\n  Idle%s         %12d     Elapsed: %12d\n",
		devChanged, string(b), idleChanged, c.Idle, c.Elapsed)
}
