// Package mix simulates the MIX computer that is described in Donald Knuth's
// "The Art of Computer Programming" (third edition).
package mix

import "log"

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

	// Less, Equal and Greater are the possible values taken by the
	// computer's comparison indicator.
	Less    = -1
	Equal   = 0
	Greater = 1

	// MemorySize is the number of memory cells in a regular MIX computer.
	MemorySize = 4000
	mBase      = MemorySize - 1
)

type Computer struct {
	*Binding
	Overflow    bool
	Comparison  int
	Elapsed     int64
	Reg         [10]Word
	Contents    []Word
	Devices     []Peripheral
	busyUntil   []int64
	m, next     int
	ctrl, trace bool
}

func NewComputer(bind *Binding) *Computer {
	if bind == nil {
		bind = DefaultBinding
	}
	return &Computer{
		Binding:   bind,
		Contents:  make([]Word, 2*MemorySize-1),
		Devices:   make([]Peripheral, 20),
		busyUntil: make([]int64, 20),
	}
}

func (c *Computer) Shutdown() error {
	var err error
	for i := range c.Devices {
		if err2 := c.unbindDevice(i); err2 != nil {
			if err == nil {
				err = err2
			} else {
				log.Println("error:", err2)
			}
		}
	}
	return err
}

func (c *Computer) validAddress(address int) bool {
	if c.ctrl {
		return address > -MemorySize && address < MemorySize
	} else {
		return address >= 0 && address < MemorySize
	}
}
