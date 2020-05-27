package mix

type Teletype struct{}

func NewTeletype() *Teletype {
	return new(Teletype)
}

func (*Teletype) Name() string {
	return "Teletype"
}

func (*Teletype) BlockSize() int {
	return 14
}

func (t *Teletype) Read([]Word) error {
	return nil
}

func (t *Teletype) Write([]Word) error {
	return nil
}

func (t *Teletype) Control(op int) error {
	return nil
}

func (t *Teletype) BusyUntil(now int64) int64 {
	return 0
}
