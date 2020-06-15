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
	return c.busyUntil[unit] > c.elapsed
}

func (c *Computer) calcTime(unit int, t int64, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	dur := c.busyUntil[unit] - c.elapsed + 1
	if dur < 1 {
		dur = 1
	}
	c.busyUntil[unit] = c.elapsed + dur + t
	return dur, nil
}
