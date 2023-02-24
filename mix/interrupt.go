// Copyright (c) 2020 Justin Flude. All rights reserved.
// Use of this source code is governed by the COPYING.md file.
package mix

import "container/heap"

const tickRate = 1000

func (c *Computer) checkInterrupt(currentState bool) {
	if c.ctrl != currentState {
		if c.ctrl {
			c.saveRegs()
			c.next = -12
		} else {
			c.loadRegs()
		}
		return
	}
	if since := c.Elapsed - c.lastTick; since >= tickRate {
		prev := int64(c.Contents[mBase-10].Int())
		now := prev - since/tickRate
		if prev > 0 && now <= 0 {
			c.schedule(c.lastTick+prev*tickRate, -11)
		}
		c.Contents[mBase-10] = NewWord(int(now))
		c.lastTick = c.Elapsed - c.Elapsed%tickRate
	}
	if !c.ctrl && c.pending.Len() > 0 && c.pending[0].when <= c.Elapsed {
		ev := heap.Pop(&c.pending).(event)
		c.saveRegs()
		c.ctrl = true
		c.next = ev.loc
	}
}

func (c *Computer) saveRegs() {
	copy(c.Contents[mBase-9:], c.Reg[:J])
	w := c.Reg[J]
	w.SetField(Spec(0, 2), NewWord(c.next))
	ov := 0
	if c.Overflow {
		ov = 1
	}
	w.SetField(Spec(3, 3), NewWord(Spec(ov, c.Comparison+1)))
	c.Contents[mBase-1] = w
}

func (c *Computer) loadRegs() {
	copy(c.Reg[:J], c.Contents[mBase-9:])
	w := c.Contents[mBase-1]
	c.Reg[J] = w.Field(Spec(4, 5))
	c.next = w.Field(Spec(1, 2)).Int()
	ovci := w.Field(Spec(3, 3)).Int()
	c.Overflow = (ovci/8 == 1)
	c.Comparison = ovci%8 - 1
}

func (c *Computer) schedule(when int64, loc int) {
	heap.Push(&c.pending, event{when, loc})
}
