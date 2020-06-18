package mix

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidAddress = errors.New("mix: invalid address")
	ErrInvalidIndex   = errors.New("mix: invalid index")
	ErrInvalidOp      = errors.New("mix: invalid operation")
)

func (c *Computer) Cycle() (err error) {
	if !c.validAddress(c.next) {
		return ErrInvalidAddress
	}
	defer func() {
		if r := recover(); r != nil {
			if err2, ok := r.(error); ok {
				err = err2
			} else {
				panic(r)
			}
		}
		if err != nil {
			err = fmt.Errorf("%w at %04d: %s", err, c.next,
				Disassemble(c.Contents[mBase+c.next]))
			if errors.Is(err, ErrHalted) {
				c.Elapsed++
				c.next++
			}
		}
	}()
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
	if c.trace {
		c.printTrace(m, c.next)
	}
	state := c.ctrl
	var t int64
	switch {
	case op == NOP:
		t = 1
	case op == ADD:
		t = c.add(aa, i, f, op, m)
	case op == SUB:
		t = c.sub(aa, i, f, op, m)
	case op == MUL:
		t = c.mul(aa, i, f, op, m)
	case op == DIV:
		t = c.div(aa, i, f, op, m)
	case op == NUM:
		t = c.num(aa, i, f, op, m)
	case op == SLA:
		t = c.sla(aa, i, f, op, m)
	case op == MOVE:
		t = c.move(aa, i, f, op, m)
	case op >= LDA && op <= LDX:
		t = c.lda(aa, i, f, op, m)
	case op >= LDAN && op <= LDXN:
		t = c.ldan(aa, i, f, op, m)
	case op >= STA && op <= STZ:
		t = c.sta(aa, i, f, op, m)
	case op == JBUS:
		t, err = c.jbus(aa, i, f, op, m)
	case op == IOC:
		t, err = c.ioc(aa, i, f, op, m)
	case op == IN:
		t, err = c.in(aa, i, f, op, m)
	case op == OUT:
		t, err = c.out(aa, i, f, op, m)
	case op == JRED:
		t, err = c.jred(aa, i, f, op, m)
	case op == JMP:
		t = c.jmp(aa, i, f, op, m)
	case op >= JA && op <= JX:
		t = c.ja(aa, i, f, op, m)
	case op >= INCA && op <= INCX:
		t = c.inca(aa, i, f, op, m)
	case op >= CMPA && op <= CMPX:
		t = c.cmpa(aa, i, f, op, m)
	}
	if err != nil {
		return err
	}
	c.Elapsed += t
	c.next++
	c.checkInterrupt(state)
	return nil
}
