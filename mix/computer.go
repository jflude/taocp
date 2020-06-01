// Package mix simulates the MIX computer described in Donald Knuth's
// "The Art of Computer Programming, Vol. 1" (second edition).
package mix

import (
	"errors"
	"log"
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
	Devices   []Peripheral
	busyUntil []int64
	elapsed   int64
	m, next   int
	trace     bool
}

var (
	// Errors returned by the CPU.
	ErrHalted             = errors.New("halted")
	ErrInvalidAddress     = errors.New("invalid address")
	ErrInvalidIndex       = errors.New("invalid index")
	ErrInvalidInstruction = errors.New("invalid instruction")
	ErrNotImplemented     = errors.New("not implemented")
)

func NewComputer() *Computer {
	return &Computer{
		CPU:       new(CPU),
		Contents:  new(Contents),
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
