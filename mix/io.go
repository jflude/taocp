package mix

func (c *Computer) jbus(aa Word, i, f, op, m int) (int, error) {
	if err := c.bindDevice(f); err != nil {
		return 0, err
	}
	c.jump(m, c.Devices[f].BusyUntil() != 0)
	return 1, nil
}

func (c *Computer) ioc(aa Word, i, f, op, m int) (int, error) {
	if err := c.bindDevice(f); err != nil {
		return 0, err
	}
	c.waitBusy(f)
	return 1, c.Devices[f].Control(m)
}

func (c *Computer) in(aa Word, i, f, op, m int) (int, error) {
	if err := c.bindDevice(f); err != nil {
		return 0, err
	}
	n := m + c.Devices[f].BlockSize()
	if m < 0 || n >= MemorySize {
		return 0, ErrInvalidAddress
	}
	c.waitBusy(f)
	return 1, c.Devices[f].Read(c.Contents[m:n])
}

func (c *Computer) out(aa Word, i, f, op, m int) (int, error) {
	if err := c.bindDevice(f); err != nil {
		return 0, err
	}
	n := m + c.Devices[f].BlockSize()
	if m < 0 || n >= MemorySize {
		return 0, ErrInvalidAddress
	}
	c.waitBusy(f)
	return 1, c.Devices[f].Write(c.Contents[m:n])
}

func (c *Computer) jred(aa Word, i, f, op, m int) (int, error) {
	if err := c.bindDevice(f); err != nil {
		return 0, err
	}
	c.jump(m, c.Devices[f].BusyUntil() == 0)
	return 1, nil
}

func (c *Computer) waitBusy(unit int) {
	if until := c.Devices[unit].BusyUntil(); until > 0 {
		c.elapsed = until
	}
}
