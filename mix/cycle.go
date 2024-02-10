package mix

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidAddress    = errors.New("mix: invalid address")
	ErrInvalidIndex      = errors.New("mix: invalid index")
	ErrInvalidOp         = errors.New("mix: invalid operation")
	ErrContentsInterlock = errors.New("mix: contents interlock")
)

func (c *Computer) Cycle() (err error) {
	defer func() {
		if r := recover(); r != nil {
			if err2, ok := r.(error); ok {
				err = err2
			} else {
				panic(r)
			}
		}
		if err != nil {
			asm := "?"
			if c.validAddress(c.next) {
				asm = Disassemble(c.Contents[mBase+c.next])
			}
			err = fmt.Errorf("%w at %04d: %s",
				err, c.next, strings.TrimSpace(asm))
		}
	}()
	c.checkAddress(c.next)
	c.checkInterlock(c.next, c.next)
	aa, i, f, op := c.Contents[mBase+c.next].UnpackOp()
	if i > 6 {
		return ErrInvalidOp
	}
	m := aa.Int()
	if i > 0 {
		m += c.Reg[i].Int()
		if m < MinIndex || m > MaxIndex {
			return ErrInvalidIndex
		}
	}
	if c.Tracer != nil && c.next >= c.Trigger {
		c.printTrace(m, c.next)
	}
	state := c.ctrl
	duration, err := microcode[op](c, aa, i, f, op, m)
	if err != nil {
		return err
	}
	c.Elapsed += duration
	c.next++
	if c.Interrupts {
		c.checkInterrupt(state)
	}
	return nil
}

func (c *Computer) nop(aa Word, i, f, op, m int) (int64, error) {
	return 1, nil
}
