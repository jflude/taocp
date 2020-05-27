package mix

type CardReader struct{}

func NewCardReader() *CardReader {
	return new(CardReader)
}

func (*CardReader) Name() string {
	return "Card Reader"
}

func (*CardReader) BlockSize() int {
	return 16
}

func (r *CardReader) Read([]Word) error {
	return nil
}

func (*CardReader) Write([]Word) error {
	return ErrInvalidOperation
}

func (r *CardReader) Control(op int) error {
	return nil
}

func (r *CardReader) BusyUntil(now int64) int64 {
	return 0
}
