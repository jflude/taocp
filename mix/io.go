package mix

func (c *Computer) jbus(aa Word, i, f, op, m int) (delay int64, err error) {
	defer func() {
		err = wrapUnitError(f, err)
	}()
	if err = c.bindDevice(f); err != nil {
		return
	}
	c.jump(m, c.isBusy(f))
	return 1, nil
}

func (c *Computer) ioc(aa Word, i, f, op, m int) (delay int64, err error) {
	defer func() {
		err = wrapUnitError(f, err)
	}()
	if err = c.bindDevice(f); err != nil {
		return
	}
	t, err := c.Devices[f].Control(m)
	return c.calcTiming(f, t, err)
}

func (c *Computer) in(aa Word, i, f, op, m int) (delay int64, err error) {
	defer func() {
		err = wrapUnitError(f, err)
	}()
	if err = c.bindDevice(f); err != nil {
		return
	}
	n := m + c.Devices[f].BlockSize()
	if !c.validAddress(m) || !c.validAddress(n) {
		return 0, ErrInvalidAddress
	}
	t, err := c.Devices[f].Read(c.Contents[mBase+m : mBase+n])
	return c.calcTiming(f, t, err)
}

func (c *Computer) out(aa Word, i, f, op, m int) (delay int64, err error) {
	defer func() {
		err = wrapUnitError(f, err)
	}()
	if err = c.bindDevice(f); err != nil {
		return
	}
	n := m + c.Devices[f].BlockSize()
	if !c.validAddress(m) || !c.validAddress(n) {
		return 0, ErrInvalidAddress
	}
	t, err := c.Devices[f].Write(c.Contents[mBase+m : mBase+n])
	return c.calcTiming(f, t, err)
}

func (c *Computer) jred(aa Word, i, f, op, m int) (delay int64, err error) {
	defer func() {
		err = wrapUnitError(f, err)
	}()
	if err = c.bindDevice(f); err != nil {
		return
	}
	c.jump(m, !c.isBusy(f))
	return 1, nil
}
