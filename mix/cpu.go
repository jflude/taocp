package mix

const (
	// The CPU registers of the MIX computer.
	A = iota
	I1
	I2
	I3
	I4
	I5
	I6
	X
	J
	// Z
)

type CPU struct {
	Reg        [10]Word
	Overflow   bool
	Comparison int
}
