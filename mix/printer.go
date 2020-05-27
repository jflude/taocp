package mix

type Printer struct{}

func NewPrinter() *Printer {
	return new(Printer)
}

func (*Printer) Name() string {
	return "Printer"
}

func (*Printer) BlockSize() int {
	return 24
}

func (*Printer) Read([]Word) error {
	return ErrInvalidOperation
}

func (p *Printer) Write([]Word) error {
	return nil
}

func (p *Printer) Control(op int) error {
	return nil
}

func (p *Printer) BusyUntil(now int64) int64 {
	return 0
}
