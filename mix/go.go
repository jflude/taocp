package mix

import "errors"

var ErrNotImplemented = errors.New("mix: not implemented")

// GoButton starts the MIX computer, as described in Ex. 26, Section 1.3.1.
// The machine can be bootstrapped from the card reader (unit 16) or the
// paper tape reader (unit 20).
func (c *Computer) GoButton(unit int) error {
	if unit != 16 && unit != 20 {
		return ErrNotImplemented
	}
	c.zeroContents()
	if _, err := c.in(0, 0, unit, IN, 0); err != nil {
		return err
	}
	c.Reg[J] = 0
	c.Elapsed = 0
	c.lastTick = 0
	c.pending = nil
	c.busyUntil[unit] = 0
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
