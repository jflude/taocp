package mix

import "errors"

type Peripheral interface {
	Name() string
	BlockSize() int
	Read(block []Word) (timing int64, err error)
	Write(block []Word) (timing int64, err error)
	Control(m int) (timing int64, err error)
	Close() error
}

// ErrInvalidCommand is returned when an I/O device is asked to do something
// it does not support, eg. requesting input from the printer.
var ErrInvalidCommand = errors.New("mix: invalid I/O command")

func (c *Computer) isBusy(unit int) bool {
	return c.busyUntil[unit] > c.Elapsed
}

func (c *Computer) calcTiming(unit int, t int64, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	delay := c.busyUntil[unit] - c.Elapsed
	if delay < 0 {
		delay = 0
	}
	c.Idle += delay
	delay++
	c.busyUntil[unit] = c.Elapsed + delay + t
	return delay, nil
}
