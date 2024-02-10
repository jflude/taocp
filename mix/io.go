package mix

func (c *Computer) jbus(aa Word, i, f, op, m int) (duration int64, err error) {
	defer func() {
		err = wrapUnitError(f, err)
	}()
	if err = c.bindDevice(f); err != nil {
		return
	}
	c.jump(m, c.isBusy(f))
	return 1, nil
}

func (c *Computer) ioc(aa Word, i, f, op, m int) (duration int64, err error) {
	defer func() {
		err = wrapUnitError(f, err)
	}()
	if err = c.bindDevice(f); err != nil {
		return
	}
	c.waitIfBusy(f)
	dur, err := c.Devices[f].Control(m)
	return 1, c.interlock(f, dur, 0, 0, err)
}

func (c *Computer) in(aa Word, i, f, op, m int) (duration int64, err error) {
	defer func() {
		err = wrapUnitError(f, err)
	}()
	if err = c.bindDevice(f); err != nil {
		return
	}
	c.waitIfBusy(f)
	n := m + c.Devices[f].BlockSize()
	c.checkAddress(m)
	c.checkAddress(n)
	c.checkInterlock(m, n)
	dur, err := c.Devices[f].Read(c.Contents[mBase+m : mBase+n])
	return 1, c.interlock(f, dur, mBase+m, mBase+n, err)
}

func (c *Computer) out(aa Word, i, f, op, m int) (duration int64, err error) {
	defer func() {
		err = wrapUnitError(f, err)
	}()
	if err = c.bindDevice(f); err != nil {
		return
	}
	c.waitIfBusy(f)
	n := m + c.Devices[f].BlockSize()
	c.checkAddress(m)
	c.checkAddress(n)
	c.checkInterlock(m, n)
	dur, err := c.Devices[f].Write(c.Contents[mBase+m : mBase+n])
	return 1, c.interlock(f, dur, mBase+m, mBase+n, err)
}

func (c *Computer) jred(aa Word, i, f, op, m int) (duration int64, err error) {
	defer func() {
		err = wrapUnitError(f, err)
	}()
	if err = c.bindDevice(f); err != nil {
		return
	}
	c.jump(m, !c.isBusy(f))
	return 1, nil
}
