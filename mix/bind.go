package mix

import (
	"errors"
	"fmt"
	"io"
	"os"
)

type Binding [21]interface{}

type readWriteSeekCloser interface {
	io.ReadWriteSeeker
	io.Closer
}

var (
	defBind = Binding{
		"tape00.mix",
		"tape01.mix",
		"tape02.mix",
		"tape03.mix",
		"tape04.mix",
		"tape05.mix",
		"tape06.mix",
		"tape07.mix",
		"drum08.mix",
		"drum09.mix",
		"drum10.mix",
		"drum11.mix",
		"disc12.mix",
		"disc13.mix",
		"disc14.mix",
		"disc15.mix",
		"reader.mix",
		"punch.mix",
		"printer.mix",
		nil,
		"paper.mix",
	}
	DefaultBinding   = &defBind
	ErrInvalidDevice = errors.New("mix: invalid I/O device")
	ErrNoDevice      = errors.New("mix: no I/O device")
)

func (c *Computer) bindDevice(unit int) error {
	if unit < 0 || unit > 20 {
		return ErrInvalidDevice
	}
	if c.Devices[unit] != nil {
		return nil
	}
	if c.bind == nil {
		c.bind = DefaultBinding
	}
	var err error
	backing := c.bind[unit]
	if file, ok := backing.(string); ok {
		var flags int
		switch unit {
		case 16:
			flags = os.O_RDONLY
		case 17, 18:
			flags = os.O_WRONLY | os.O_APPEND
		default:
			flags = os.O_RDWR
		}
		backing, err = os.OpenFile(file, os.O_CREATE|flags, 0644)
		if err != nil {
			return err
		}
	} else if unit == 19 && backing == nil {
		backing = console{}
	}
	if backing == nil {
		return fmt.Errorf("%w: unit %d", ErrNoDevice, unit)
	}
	var p Peripheral
	switch {
	case unit >= 0 && unit <= 7:
		p, err = NewTape(backing.(readWriteSeekCloser), unit)
	case unit >= 8 && unit <= 11:
		p, err = NewDrum(backing.(readWriteSeekCloser), unit, c)
	case unit >= 12 && unit <= 15:
		p, err = NewDisc(backing.(readWriteSeekCloser), unit, c)
	case unit == 16:
		p, err = NewCardReader(backing.(io.ReadCloser))
	case unit == 17:
		p, err = NewCardPunch(backing.(io.WriteCloser))
	case unit == 18:
		p, err = NewPrinter(backing.(io.WriteCloser))
	case unit == 19:
		p, err = NewTeletype(backing.(io.ReadWriteCloser))
	case unit == 20:
		p, err = NewPaperTape(backing.(readWriteSeekCloser))
	}
	if err == nil {
		c.Devices[unit] = p
	}
	return err
}

func (c *Computer) unbindDevice(unit int) error {
	if unit < 0 || unit > 20 {
		return ErrInvalidDevice
	}
	d := c.Devices[unit]
	if d == nil {
		return nil
	}
	c.Devices[unit] = nil
	return d.Close()
}
