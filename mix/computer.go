// Package mix simulates the MIX computer described in Donald Knuth's
// "The Art of Computer Programming, Vol. 1" (second edition).
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

	// MemorySize is the number of memory cells in a MIX computer.
	MemorySize = 4000
)

type CPU struct {
	Reg        [10]Word
	Overflow   bool
	Comparison int
}

type Contents [MemorySize]Word

type Computer struct {
	*CPU
	*Contents
	*Binding
	Devices   []Peripheral
	busyUntil []int64
	elapsed   int64
	m, next   int
	trace     bool
}

func NewComputer(bind *Binding) *Computer {
	if bind == nil {
		bind = DefaultBinding
	}
	return &Computer{
		CPU:       new(CPU),
		Contents:  new(Contents),
		Binding:   bind,
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
