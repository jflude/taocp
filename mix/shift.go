// Copyright (c) 2020 Justin Flude. All rights reserved.
// Use of this source code is governed by the COPYING.md file.
package mix

func (c *Computer) sla(aa Word, i, f, op, m int) int64 {
	if m < 0 {
		panic(ErrInvalidOp)
	}
	switch f {
	case 0: // SLA
		var w Word
		ShiftBitsLeft(&c.Reg[A], &w, 6*m)
	case 1: // SRA
		var w Word
		ShiftBitsRight(&c.Reg[A], &w, 6*m)
	case 2: // SLAX
		ShiftBitsLeft(&c.Reg[A], &c.Reg[X], 6*m)
	case 3: // SRAX
		ShiftBitsRight(&c.Reg[A], &c.Reg[X], 6*m)
	case 4: // SLC
		RotateBitsLeft(&c.Reg[A], &c.Reg[X], 6*m)
	case 5: // SRC
		RotateBitsRight(&c.Reg[A], &c.Reg[X], 6*m)
	case 6: // SLB (see Section 4.5.2)
		ShiftBitsLeft(&c.Reg[A], &c.Reg[X], m)
	case 7: // SRB (see Section 4.5.2)
		ShiftBitsRight(&c.Reg[A], &c.Reg[X], m)
	default:
		panic(ErrInvalidOp)
	}
	return 2
}
