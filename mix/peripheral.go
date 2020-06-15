package mix

import "errors"

type Peripheral interface {
	Name() string
	BlockSize() int
	Read(block []Word) (int64, error)
	Write(block []Word) (int64, error)
	Control(m int) (int64, error)
	Close() error
}

// ErrInvalidCommand is returned when an I/O device is asked to do something
// it does not support, eg. requesting input from the printer.
var ErrInvalidCommand = errors.New("mix: invalid I/O command")

func (c *Computer) isBusy(unit int) bool {
	return c.busyUntil[unit] > c.Elapsed
}

// TODO: check timings of all I/O operations to see if typical for ~1970
func (c *Computer) calcTiming(unit int, t int64, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	delay := c.busyUntil[unit] - c.Elapsed + 1
	if delay < 1 {
		delay = 1
	}
	c.busyUntil[unit] = c.Elapsed + delay + t
	return delay, nil
}
