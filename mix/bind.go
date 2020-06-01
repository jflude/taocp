package mix

func (c *Computer) bindDevice(unit int) error {
	if unit < 0 || unit >= 20 {
		return ErrInvalidDevice
	}
	if c.Devices[unit] != nil {
		return nil
	}
	var p Peripheral
	var err error
	switch {
	case unit >= 0 && unit <= 7:
		p, err = NewTape(nil, unit)
	case unit >= 8 && unit <= 15:
		p, err = NewDrum(nil, unit, c)
	case unit == 16:
		p, err = NewCardReader(nil)
	case unit == 17:
		p, err = NewCardPunch(nil)
	case unit == 18:
		p, err = NewPrinter(nil)
	case unit == 19:
		p, err = NewTeletype(nil)
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
