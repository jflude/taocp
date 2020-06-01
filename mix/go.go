package mix

import "errors"

// GoButton starts the MIX computer, as described in Ex. 26, section 1.3.1.
func (c *Computer) GoButton() error {
	for {
		if err := c.Cycle(); err != nil {
			if errors.Is(err, ErrHalted) {
				return nil
			}
			return err
		}
	}
}
