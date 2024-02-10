package mix

// GoButton starts the MIX computer, as described in Ex. 26, Section 1.3.1.
// The machine can be bootstrapped only from the card reader or paper tape.
func (c *Computer) GoButton() error {
	if c.halted {
		c.halted = false
		c.next++
		return c.resume()
	}
	if c.BootFrom != CardReaderUnit && c.BootFrom != PaperTapeUnit {
		return ErrInvalidUnit
	}
	for i := range c.Reg {
		c.Reg[i] = 0
	}
	c.zeroContents()
	c.unlockContents()
	if _, err := c.in(0, 0, c.BootFrom, IN, 0); err != nil {
		return err
	}
	c.unlockContents()
	c.ctrl = c.Interrupts
	c.Reg[J] = 0
	c.Elapsed = 0
	c.lastTick = 0
	c.pending = nil
	for i := range c.busyUntil {
		c.busyUntil[i] = 0
	}
	c.next = 0
	return c.resume()
}

func (c *Computer) resume() error {
	for {
		if err := c.Cycle(); err != nil {
			return err
		}
	}
}
