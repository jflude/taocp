package mix

func (c *Computer) jbus(aa Word, i, f, op, m int) int {
	if f >= len(c.Devices) {
		panic(ErrInvalidInstruction)
	}
	c.jump(m, c.Devices[f].BusyUntil(c.elapsed) != 0)
	return 1
}

func (c *Computer) ioc(aa Word, i, f, op, m int) int {
	if f >= len(c.Devices) {
		panic(ErrInvalidInstruction)
	}
	c.waitBusy(f)
	c.Devices[f].Control(m)
	return 1
}

func (c *Computer) in(aa Word, i, f, op, m int) int {
	if f >= len(c.Devices) {
		panic(ErrInvalidInstruction)
	}
	n := m + c.Devices[f].BlockSize()
	if m < 0 || n >= MemorySize {
		panic(ErrInvalidAddress)
	}
	c.waitBusy(f)
	c.Devices[f].Read(c.Contents[m:n])
	return 1
}

func (c *Computer) out(aa Word, i, f, op, m int) int {
	if f >= len(c.Devices) {
		panic(ErrInvalidInstruction)
	}
	n := m + c.Devices[f].BlockSize()
	if m < 0 || n >= MemorySize {
		panic(ErrInvalidAddress)
	}
	c.waitBusy(f)
	c.Devices[f].Write(c.Contents[m:n])
	return 1
}

func (c *Computer) jred(aa Word, i, f, op, m int) int {
	if f >= len(c.Devices) {
		panic(ErrInvalidInstruction)
	}
	c.jump(m, c.Devices[f].BusyUntil(c.elapsed) == 0)
	return 1
}

func (c *Computer) waitBusy(unit int) {
	if until := c.Devices[unit].BusyUntil(c.elapsed); until > 0 {
		c.elapsed = until
	}
}
