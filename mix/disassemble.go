package mix

import "fmt"

var (
	regs   = []string{"A", "1", "2", "3", "4", "5", "6", "X", "J", "Z"}
	ariths = []string{"ADD", "SUB", "MUL", "DIV"}
	nums   = []string{"NUM", "CHAR", "HLT", "AND", "OR", "XOR", "FLOT",
		"FIX", "", "INT"}
	shifts = []string{"SLA", "SRA", "SLAX", "SRAX", "SLC", "SRC", "SLB",
		"SRB"}
	ios   = []string{"JBUS", "IOC", "IN", "OUT", "JRED"}
	jumps = []string{"JMP", "JSJ", "JOV", "JNV", "JL", "JE", "JG", "JGE",
		"JNE", "JLE"}
	conds = []string{"N", "Z", "P", "NN", "NE", "NP", "E", "O"}
	incs  = []string{"INC", "DEC", "ENT", "ENN"}
)

func Disassemble(w Word) string {
	aa, i, f, op := w.UnpackOp()
	switch {
	case op == NOP:
		return noArg("NOP")
	case op >= ADD && op <= DIV:
		s := ariths[op-ADD]
		if f == 6 {
			s = "F" + s
		}
		return hasSpec(s, aa, i, f)
	case op == NUM:
		if f < len(nums) && nums[f] != "" {
			return noArg(nums[f])
		}
	case op == SLA:
		if f < len(shifts) {
			return noField(shifts[f], aa, i)
		}
	case op == MOVE:
		return hasField("MOVE", aa, i, f)
	case op >= LDA && op <= LDX:
		return hasSpec("LD"+regs[op-LDA], aa, i, f)
	case op >= LDAN && op <= LDXN:
		return hasSpec("LD"+regs[op-LDAN]+"N", aa, i, f)
	case op >= STA && op <= STZ:
		var usual int
		if op == STJ {
			usual = 2
		} else {
			usual = 5
		}
		return hasSpecUsual("ST"+regs[op-STA], aa, i, f, usual)
	case op >= JBUS && op <= JRED:
		return hasField(ios[op-JBUS], aa, i, f)
	case op == JMP:
		if f < len(jumps) {
			return noField(jumps[f], aa, i)
		}
	case op >= JA && op <= JX:
		if f < len(conds) {
			return noField("J"+regs[op-JA]+conds[f], aa, i)
		}
	case op >= INCA || op <= INCX:
		if f < len(incs) {
			return noField(incs[f]+regs[op-INCA], aa, i)
		}
	case op >= CMPA && op <= CMPX:
		var s string
		if op == CMPA && f == 6 {
			s = "FCMP"
		} else {
			s = "CMP" + regs[op-CMPA]
		}
		return hasSpec(s, aa, i, f)
	}
	return fmt.Sprintf("CON  %v", w)
}

func noArg(op string) string {
	return fmt.Sprintf("%-4s", op)
}

func noField(op string, aa Word, i int) string {
	return fmt.Sprintf("%-4s %d%s", op, aa.Int(), hasIndex(i))
}

func hasField(op string, aa Word, i, f int) string {
	return fmt.Sprintf("%-4s %d%s(%d)", op, aa.Int(), hasIndex(i), f)
}

func hasSpec(op string, aa Word, i, f int) string {
	return hasSpecUsual(op, aa, i, f, 5)
}

func hasSpecUsual(op string, aa Word, i, f, usual int) string {
	var fld string
	if f != usual {
		fld = fmt.Sprintf("(%d:%d)", f/8, f%8)
	}
	return fmt.Sprintf("%-4s %d%s%s", op, aa.Int(), hasIndex(i), fld)
}

func hasIndex(i int) string {
	if i != 0 {
		return fmt.Sprintf(",%d", i)
	}
	return ""
}
