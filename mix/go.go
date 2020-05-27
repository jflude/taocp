package mix

// GoButton starts the MIX computer, as described in exercise 26, section 1.3.1.
func (c *Computer) GoButton() error {
	for {
		if err := c.Cycle(); err != nil {
			if err == ErrHalted {
				return nil
			}
			return err
		}
	}
}
