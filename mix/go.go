package mix

import "errors"

// GoButton starts the MIX computer, as described in Ex. 26, Section 1.3.1.
// The machine can be only be bootstrapped from a card (unit 16) or a paper
// tape (unit 19).
func (c *Computer) GoButton(unit int) error {
	if unit != 16 && unit != 19 {
		return ErrNotImplemented
	}
	if _, err := c.in(0, 0, unit, IN, 0); err != nil {
		return err
	}
	c.Reg[J] = 0
	c.next = 0
	c.elapsed = 0
	c.busyUntil[unit] = 0
	return c.resume()
}

func (c *Computer) resume() error {
	for {
		if err := c.Cycle(); err != nil {
			if errors.Is(err, ErrHalted) {
				return nil
			}
			return err
		}
	}
}
