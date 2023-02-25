// Copyright (c) 2020 Justin Flude. All rights reserved.
// Use of this source code is governed by the COPYING.md file.
package mix

import "errors"

type Peripheral interface {
	Name() string
	BlockSize() int
	Read(block []Word) (duration int64, err error)
	Write(block []Word) (duration int64, err error)
	Control(m int) (duration int64, err error)
	Close() error
}

// ErrInvalidCommand is returned when an I/O device is asked to do something
// it does not support, eg. requesting input from the printer.
var ErrInvalidCommand = errors.New("mix: invalid I/O command")

func (c *Computer) isBusy(unit int) bool {
	return c.busyUntil[unit] > c.Elapsed
}

func (c *Computer) waitIfBusy(unit int) {
	if wait := c.busyUntil[unit] - c.Elapsed; wait > 0 {
		c.Idle += wait
		c.Elapsed = c.busyUntil[unit]
	}
}

func (c *Computer) interlock(unit int, dur int64, m, n int, err error) error {
	if err != nil {
		return err
	}
	c.busyUntil[unit] = c.Elapsed + dur + 1
	c.lockContents(m, n, c.busyUntil[unit])
	if c.Interrupts && c.ctrl {
		c.schedule(c.busyUntil[unit], -20-unit)
	}
	return nil
}
