package mix

import "errors"

var ErrNotImplemented = errors.New("mix: not implemented")

func (c *Computer) Save(filename string) error {
	return ErrNotImplemented // TODO
}

func Load(filename string) (*Computer, error) {
	return nil, ErrNotImplemented // TODO
}
