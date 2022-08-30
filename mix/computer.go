// Copyright (c) 2020 Justin Flude. All rights reserved.
// Use of this source code is governed by the COPYING.md file.
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

	// DeviceCount is the number of I/O devices support by the computer.
	DeviceCount = 21
)

type Computer struct {
	Overflow      bool
	Comparison    int
	Elapsed, Idle int64
	Reg           [10]Word
	Contents      []Word
	Devices       []Peripheral
	busyUntil     []int64
	bind          *Binding
	m, next       int
	ctrl, trace   bool
	Interrupts    bool
	lastTick      int64
	pending       priority
}

func NewComputer() *Computer {
	return &Computer{
		Contents:  make([]Word, 2*MemorySize-1),
		Devices:   make([]Peripheral, DeviceCount),
		busyUntil: make([]int64, DeviceCount),
	}
}

func (c *Computer) Bind(b *Binding) error {
	if err := c.unbindAll(); err != nil {
		return err
	}
	c.bind = b
	return nil
}

func (c *Computer) unbindAll() error {
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

func (c *Computer) Shutdown() error {
	return c.unbindAll()
}

func (c *Computer) validAddress(address int) bool {
	if c.ctrl {
		return address > -MemorySize && address < MemorySize
	} else {
		return address >= 0 && address < MemorySize
	}
}

func (c *Computer) zeroContents() {
	for i := range c.Contents {
		c.Contents[i] = 0
	}
}
