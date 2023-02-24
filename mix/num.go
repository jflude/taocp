// Copyright (c) 2020 Justin Flude. All rights reserved.
// Use of this source code is governed by the COPYING.md file.
package mix

import (
	"errors"
	"strconv"
)

var (
	ErrHalted       = errors.New("mix: halted")
	ErrNoInterrupts = errors.New("mix: interrupts are disabled")
)

func (c *Computer) num(aa Word, i, f, op, m int) int64 {
	switch f {
	case 0: // NUM
		var a, x int64
		for i := 1; i <= 5; i++ {
			f := Spec(i, i)
			a = 10*a + int64(c.Reg[A].Field(f).Int())%10
			x = 10*x + int64(c.Reg[X].Field(f).Int())%10
		}
		a = a*100000 + x
		if a > MaxWord {
			a %= MaxWord + 1
			c.Overflow = true
		}
		c.Reg[A].SetField(Spec(1, 5), NewWord(int(a)))
		return 10
	case 1: // CHAR
		v := strconv.Itoa(c.Reg[A].Field(Spec(1, 5)).Int())
		if l := len(v); l < 10 {
			v = "000000000"[:10-l] + v
		}
		for i := 0; i < 5; i++ {
			f := Spec(i+1, i+1)
			a := NewWord(utf2mix[rune(v[i])])
			x := NewWord(utf2mix[rune(v[i+5])])
			c.Reg[A].SetField(f, a)
			c.Reg[X].SetField(f, x)
		}
		return 10
	case 2: // HLT
		c.Elapsed += 10
		now := c.Elapsed
		for i := range c.Devices { // finish any I/O operations
			if c.busyUntil[i] > c.Elapsed {
				c.Elapsed = c.busyUntil[i]
			}
		}
		c.Idle += c.Elapsed - now
		c.halted = true
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
		c.callWithOvCheck1(FixedToFloat)
		return 3
	case 7: // FIX
		c.callWithOvCheck1(FloatToFixed)
		return 3
	case 9: // INT
		if !c.Interrupts { // see Ex. 18, Section 1.4.4
			panic(ErrNoInterrupts)
		}
		c.ctrl = !c.ctrl
		return 2
	default:
		panic(ErrInvalidOp)
	}
}

func (c *Computer) callWithOvCheck1(f func(Word) (Word, bool)) {
	var ov bool
	c.Reg[A], ov = f(c.Reg[A])
	c.Overflow = c.Overflow || ov
}
