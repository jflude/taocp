package mix

import "fmt"

type Tape struct {
	name string
}

func NewTape(unit int) *Tape {
	return &Tape{name: fmt.Sprintf("Tape #%d", unit)}
}

func (t *Tape) Name() string {
	return t.name
}

func (*Tape) BlockSize() int {
	return 100
}

func (t *Tape) Read([]Word) error {
	return nil
}

func (t *Tape) Write([]Word) error {
	return nil
}

func (t *Tape) Control(op int) error {
	return nil
}

func (t *Tape) BusyUntil(now int64) int64 {
	return 0
}
