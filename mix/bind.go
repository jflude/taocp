// Copyright (c) 2020 Justin Flude. All rights reserved.
// Use of this source code is governed by the COPYING.md file.
package mix

import (
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	// The I/O devices supported by the MIX 1009 computer.
	Tape0Unit = iota
	Tape1Unit
	Tape2Unit
	Tape3Unit
	Tape4Unit
	Tape5Unit
	Tape6Unit
	Tape7Unit
	Drum8Unit
	Drum9Unit
	Drum10Unit
	Drum11Unit
	Disc12Unit
	Disc13Unit
	Disc14Unit
	Disc15Unit
	CardReaderUnit
	CardPunchUnit
	PrinterUnit
	TeletypeUnit
	PaperTapeUnit
	DeviceCount
)

type Binding [DeviceCount]interface{}

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
	DefaultBinding = &defBind
	ErrInvalidUnit = errors.New("mix: invalid I/O unit")
	ErrNoDevice    = errors.New("mix: no I/O device")
)

func (c *Computer) bindDevice(unit int) error {
	if unit < Tape0Unit || unit > PaperTapeUnit {
		return ErrInvalidUnit
	}
	if c.Devices[unit] != nil {
		return nil
	}
	if c.bind == nil {
		c.bind = DefaultBinding
	}
	var err error
	backing := c.bind[unit]
	if unit == TeletypeUnit && backing == nil {
		backing = console{}
	} else if file, ok := backing.(string); ok {
		var flags int
		switch unit {
		case CardReaderUnit:
			flags = os.O_RDONLY
		case CardPunchUnit:
			flags = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
		case PrinterUnit:
			flags = os.O_WRONLY | os.O_CREATE | os.O_APPEND
		default:
			flags = os.O_RDWR | os.O_CREATE
		}
		if backing, err = os.OpenFile(file, flags, 0644); err != nil {
			return err
		}
	} else {
		return ErrNoDevice
	}
	var p Peripheral
	switch {
	case unit >= Tape0Unit && unit <= Tape7Unit:
		p, err = NewTape(backing.(readWriteSeekCloser), unit)
	case unit >= Drum8Unit && unit <= Drum11Unit:
		p, err = NewDrum(backing.(readWriteSeekCloser), unit, c)
	case unit >= Disc12Unit && unit <= Disc15Unit:
		p, err = NewDisc(backing.(readWriteSeekCloser), unit, c)
	case unit == CardReaderUnit:
		p, err = NewCardReader(backing.(io.Reader))
	case unit == CardPunchUnit:
		p, err = NewCardPunch(backing.(io.WriteCloser))
	case unit == PrinterUnit:
		p, err = NewPrinter(backing.(io.WriteCloser))
	case unit == TeletypeUnit:
		p, err = NewTeletype(backing.(io.ReadWriteCloser))
	case unit == PaperTapeUnit:
		p, err = NewPaperTape(backing.(readWriteSeekCloser))
	}
	if err == nil {
		c.Devices[unit] = p
	}
	return err
}

func (c *Computer) unbindDevice(unit int) (err error) {
	defer func() {
		err = wrapUnitError(unit, err)
	}()
	if unit < Tape0Unit || unit > PaperTapeUnit {
		err = ErrInvalidUnit
		return
	}
	if c.Devices[unit] != nil {
		d := c.Devices[unit]
		c.Devices[unit] = nil
		err = d.Close()
	}
	return
}

func wrapUnitError(unit int, err error) error {
	if err != nil {
		err = fmt.Errorf("unit %d: %w", unit, err)
	}
	return err
}
