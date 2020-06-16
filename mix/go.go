package mix

import "errors"

var ErrNotImplemented = errors.New("mix: not implemented")

// GoButton starts the MIX computer, as described in Ex. 26, Section 1.3.1.
// The machine can only be bootstrapped from the card reader (unit 16) or
// the paper tape reader (unit 20).
func (c *Computer) GoButton(unit int) error {
	if unit != 16 && unit != 20 {
		return ErrNotImplemented
	}
	for i := 0; i < len(c.Contents); i++ {
		c.Contents[i] = 0
	}
	if _, err := c.in(0, 0, unit, IN, 0); err != nil {
		return err
	}
	c.Reg[J] = 0
	c.next = 0
	c.Elapsed = 0
	c.busyUntil[unit] = 0
	return c.resume()
}

func (c *Computer) resume() error {
	for {
		if err := c.Cycle(); err != nil {
			return err
		}
	}
}
