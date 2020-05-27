package mix

import "fmt"

type Disc struct {
	name string
}

func NewDisc(unit int) *Disc {
	return &Disc{name: fmt.Sprintf("Disc #%d", unit)}
}

func (d *Disc) Name() string {
	return d.name
}

func (*Disc) BlockSize() int {
	return 100
}

func (d *Disc) Read([]Word) error {
	return nil
}

func (d *Disc) Write([]Word) error {
	return nil
}

func (d *Disc) Control(op int) error {
	return nil
}

func (d *Disc) BusyUntil(now int64) int64 {
	return 0
}
