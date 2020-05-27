package mix

type CardPunch struct{}

func NewCardPunch() *CardPunch {
	return new(CardPunch)
}

func (*CardPunch) Name() string {
	return "Card Punch"
}

func (*CardPunch) BlockSize() int {
	return 16
}

func (*CardPunch) Read([]Word) error {
	return ErrInvalidOperation
}

func (p *CardPunch) Write([]Word) error {
	return nil
}

func (p *CardPunch) Control(op int) error {
	return nil
}

func (p *CardPunch) BusyUntil(now int64) int64 {
	return 0
}
