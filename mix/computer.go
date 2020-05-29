// Package mix simulates the MIX computer described in Donald Knuth's
// "The Art of Computer Programming, Vol. 1" (second edition).
package mix

import (
	"errors"
	"fmt"
)

const (
	// MemorySize is the number of memory cells in a MIX computer.
	MemorySize = 4000

	// Less, Equal and Greater are the possible values taken by the
	// computer's comparison indicator.
	Less    = -1
	Equal   = 0
	Greater = 1
)

const (
	// The CPU registers of the MIX computer.
	A = iota
	I1
	I2
	I3
	I4
	I5
	I6
	X
	J
	// Z
)

type CPU struct {
	Reg        [10]Word
	Overflow   bool
	Comparison int
}

type Contents [MemorySize]Word

type Peripheral interface {
	Name() string
	BlockSize() int
	Read([]Word) error
	Write([]Word) error
	Control(op int) error
	BusyUntil(now int64) (until int64)
}

type Computer struct {
	*CPU
	*Contents
	Devices []Peripheral
	next    int
	elapsed int64
}

var (
	// Errors returned by the CPU.
	ErrHalted             = errors.New("halted")
	ErrInvalidAddress     = errors.New("invalid address")
	ErrInvalidIndex       = errors.New("invalid index")
	ErrInvalidInstruction = errors.New("invalid instruction")
	ErrNotImplemented     = errors.New("not implemented")

	// Errors returned by peripheral devices.
	ErrInvalidControl   = errors.New("invalid I/O control")
	ErrInvalidOperation = errors.New("invalid I/O operation")
)

func NewComputer() *Computer {
	d := make([]Peripheral, 20)
	for i := 0; i < 8; i++ {
		d[i] = NewTape(i)
	}
	for i := 8; i < 16; i++ {
		d[i] = NewDisc(i)
	}
	d[16] = NewCardReader()
	d[17] = NewCardPunch()
	d[18] = NewPrinter()
	d[19] = NewTeletype()
	return &Computer{CPU: new(CPU), Contents: new(Contents), Devices: d}
}

func (c *Computer) String() string {
	s := fmt.Sprintf(" A: %10v (%#v)   X: %10v (%#v)\n",
		c.Reg[A], c.Reg[A], c.Reg[X], c.Reg[X])
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
	s += fmt.Sprintf("I1:       %4v (%#v)     OV: %s  CI: %s\n",
		c.Reg[I1], c.Reg[I1], ov, ci)
	n := c.next
	for i := 2; i <= 6; i, n = i+1, n+1 {
		if n < MemorySize {
			s += fmt.Sprintf("I%d:       %4v (%#v)   %4d: %#v\n",
				i, c.Reg[i], c.Reg[i], n, c.Contents[n])
		} else {
			s += fmt.Sprintf("I%d:       %4v (%#v)\n",
				i, c.Reg[i], c.Reg[i])
		}
	}
	b := make([]byte, len(c.Devices))
	for i := 0; i < 20; i++ {
		if c.Devices[i].BusyUntil(c.elapsed) != 0 {
			b[i] += 'B'
		} else {
			b[i] = '.'
		}
	}
	s += fmt.Sprintf("Elapsed: %d    Device: %s\n\n",
		c.elapsed, string(b))
	return s
}
