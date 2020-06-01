package mix

func (c *Computer) bindDevice(f int) error {
	if f < 0 || f >= 20 {
		return ErrInvalidDevice
	}
	if c.Devices[f] != nil {
		return nil
	}
	var p Peripheral
	var err error
	switch {
	case f >= 0 && f <= 7:
		p, err = NewTape(c, nil, f)
	case f >= 8 && f <= 15:
		p, err = NewDrum(c, nil, f)
	case f == 16:
		p, err = NewCardReader(c, nil)
	case f == 17:
		p, err = NewCardPunch(c, nil)
	case f == 18:
		p, err = NewPrinter(c, nil)
	case f == 19:
		p, err = NewTeletype(c, nil)
	}
	if err == nil {
		c.Devices[f] = p
	}
	return err
}

func (c *Computer) unbindDevice(f int) error {
	if f < 0 || f >= 20 {
		return ErrInvalidDevice
	}
	d := c.Devices[f]
	if d == nil {
		return nil
	}
	c.Devices[f] = nil
	return d.Close()
}
