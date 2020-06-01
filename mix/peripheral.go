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

var (
	// Errors returned by peripheral devices.
	ErrInvalidDevice    = errors.New("invalid I/O device")
	ErrInvalidControl   = errors.New("invalid I/O control")
	ErrInvalidOperation = errors.New("invalid I/O operation")
)

func (c *Computer) isBusy(unit int) bool {
	return c.busyUntil[unit] > c.elapsed
}

func (c *Computer) calcTime(unit int, t int64, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	var dur int64 = 1
	if c.busyUntil[unit] > c.elapsed {
		dur += c.busyUntil[unit] - c.elapsed
	}
	c.busyUntil[unit] = c.elapsed + dur + t
	return dur, nil
}
