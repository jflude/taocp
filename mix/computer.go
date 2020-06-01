// Package mix simulates the MIX computer described in Donald Knuth's
// "The Art of Computer Programming, Vol. 1" (second edition).
package mix

import (
	"errors"
	"log"
)

const (
	// Less, Equal and Greater are the possible values taken by the
	// computer's comparison indicator.
	Less    = -1
	Equal   = 0
	Greater = 1
)

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
