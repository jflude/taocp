package mix

import "errors"

type Peripheral interface {
	Name() string
	BlockSize() int
	Read(block []Word) error
	Write(block []Word) error
	Control(m int) error
	BusyUntil() int64
	Close() error
}

var (
	// Errors returned by peripheral devices.
	ErrInvalidDevice    = errors.New("invalid I/O device")
	ErrInvalidControl   = errors.New("invalid I/O control")
	ErrInvalidOperation = errors.New("invalid I/O operation")
)
