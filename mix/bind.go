package mix

import "errors"

type Binding [20]string

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
		"drum12.mix",
		"drum13.mix",
		"drum14.mix",
		"drum15.mix",
		"reader.mix",
		"punch.mix",
		"printer.mix",
		"",
	}
	DefaultBinding   = &defBind
	ErrInvalidDevice = errors.New("mix: invalid I/O device")
)

func (c *Computer) bindDevice(unit int) error {
	if unit < 0 || unit >= 20 {
		return ErrInvalidDevice
	}
	if c.Devices[unit] != nil {
		return nil
	}
	f := c.Binding[unit]
	var p Peripheral
	var err error
	switch {
	case unit >= 0 && unit <= 7:
		p, err = NewTape(f, unit)
	case unit >= 8 && unit <= 15:
		p, err = NewDrum(f, unit, c)
	case unit == 16:
		p, err = NewCardReader(f)
	case unit == 17:
		p, err = NewCardPunch(f)
	case unit == 18:
		p, err = NewPrinter(f)
	case unit == 19:
		p, err = NewTeletype(f)
	}
	if err == nil {
		c.Devices[unit] = p
	}
	return err
}

func (c *Computer) unbindDevice(unit int) error {
	if unit < 0 || unit >= 20 {
		return ErrInvalidDevice
	}
	d := c.Devices[unit]
	if d == nil {
		return nil
	}
	c.Devices[unit] = nil
	return d.Close()
}
