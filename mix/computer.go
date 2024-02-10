
// Package mix simulates the MIX 1009 computer as described in
// Donald Knuth's "The Art of Computer Programming" (third edition).
package mix

import (
	"io"
	"log"
)

const (
	// The registers of the MIX 1009 CPU.
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
	Overflow      bool
	Comparison    int
	Elapsed, Idle int64
	Reg           [10]Word
	Contents      []Word
	attributes    []attribute
	Devices       []Peripheral
	busyUntil     []int64
	bind          *Binding
	m, next       int
	ctrl, halted  bool
	Tracer        io.WriteCloser
	Trigger       int
	lastDevMask   uint
	lastIdle      int64
	Interrupts    bool
	lastTick      int64
	pending       priority
	BootFrom      int
}

type attribute struct {
	lockUntil int64
}

func NewComputer() *Computer {
	return &Computer{
		Contents:   make([]Word, 2*MemorySize-1),
		attributes: make([]attribute, 2*MemorySize-1),
		Devices:    make([]Peripheral, DeviceCount),
		busyUntil:  make([]int64, DeviceCount),
		BootFrom:   CardReaderUnit,
		Trigger:    32,
	}
}

func (c *Computer) Bind(b *Binding) error {
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
	c.bind = b
	return err
}

func (c *Computer) Shutdown() error {
	return c.Bind(nil)
}

func (c *Computer) validAddress(address int) bool {
	if c.ctrl {
		return address > -MemorySize && address < MemorySize
	} else {
		return address >= 0 && address < MemorySize
	}
}

func (c *Computer) checkAddress(m int) {
	if !c.validAddress(m) {
		panic(ErrInvalidAddress)
	}
}

func (c *Computer) zeroContents() {
	for i := range c.Contents {
		c.Contents[i] = 0
	}
}

func (c *Computer) unlockContents() {
	for i := range c.attributes {
		c.attributes[i].lockUntil = 0
	}
}

func (c *Computer) lockContents(m, n int, until int64) {
	for i := m; i < n; i++ {
		c.attributes[i].lockUntil = until
	}
}

func (c *Computer) checkInterlock(m, n int) {
	if m == n {
		if c.attributes[mBase+m].lockUntil > c.Elapsed {
			panic(ErrContentsInterlock)
		}
	} else {
		for i := mBase + m; i < mBase+n; i++ {
			if c.attributes[i].lockUntil > c.Elapsed {
				panic(ErrContentsInterlock)
			}
		}
	}
}
