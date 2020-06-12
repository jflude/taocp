package mixal

import "github.com/jflude/gnuth/mix"

// The opcodes for the MIX computer (see Table 1, Section 1.3.1)
var opcodes = map[string]struct {
	c, f int // C-value and default F-value
}{
	"NOP":  {mix.NOP, 5},
	"ADD":  {mix.ADD, 5},
	"FADD": {mix.ADD, 6},
	"SUB":  {mix.SUB, 5},
	"FSUB": {mix.SUB, 6},
	"MUL":  {mix.MUL, 5},
	"FMUL": {mix.MUL, 6},
	"DIV":  {mix.DIV, 5},
	"FDIV": {mix.DIV, 6},
	"NUM":  {mix.NUM, 0},
	"CHAR": {mix.NUM, 1},
	"HLT":  {mix.NUM, 2},
	//"INT":  {mix.NUM, 3},
	"SLA":  {mix.SLA, 0},
	"SRA":  {mix.SLA, 1},
	"SLAX": {mix.SLA, 2},
	"SRAX": {mix.SLA, 3},
	"SLC":  {mix.SLA, 4},
	"SRC":  {mix.SLA, 5},
	"MOVE": {mix.MOVE, 0},
	"LDA":  {mix.LDA, 5},
	"LD1":  {mix.LD1, 5},
	"LD2":  {mix.LD2, 5},
	"LD3":  {mix.LD3, 5},
	"LD4":  {mix.LD4, 5},
	"LD5":  {mix.LD5, 5},
	"LD6":  {mix.LD6, 5},
	"LDX":  {mix.LDX, 5},
	"LDAN": {mix.LDAN, 5},
	"LD1N": {mix.LD1N, 5},
	"LD2N": {mix.LD2N, 5},
	"LD3N": {mix.LD3N, 5},
	"LD4N": {mix.LD4N, 5},
	"LD5N": {mix.LD5N, 5},
	"LD6N": {mix.LD6N, 5},
	"LDXN": {mix.LDXN, 5},
	"STA":  {mix.STA, 5},
	"ST1":  {mix.ST1, 5},
	"ST2":  {mix.ST2, 5},
	"ST3":  {mix.ST3, 5},
	"ST4":  {mix.ST4, 5},
	"ST5":  {mix.ST5, 5},
	"ST6":  {mix.ST6, 5},
	"STX":  {mix.STX, 5},
	"STJ":  {mix.STJ, 2},
	"STZ":  {mix.STZ, 5},
	"JBUS": {mix.JBUS, 0},
	"IOC":  {mix.IOC, 0},
	"IN":   {mix.IN, 0},
	"OUT":  {mix.OUT, 0},
	"JRED": {mix.JRED, 0},
	"JMP":  {mix.JMP, 0},
	"JSJ":  {mix.JMP, 1},
	"JOV":  {mix.JMP, 2},
	"JNOV": {mix.JMP, 3},
	"JL":   {mix.JMP, 4},
	"JE":   {mix.JMP, 5},
	"JG":   {mix.JMP, 6},
	"JGE":  {mix.JMP, 7},
	"JNE":  {mix.JMP, 8},
	"JLE":  {mix.JMP, 9},
	"JAN":  {mix.JA, 0},
	"JAZ":  {mix.JA, 1},
	"JAP":  {mix.JA, 2},
	"JANN": {mix.JA, 3},
	"JANZ": {mix.JA, 4},
	"JANP": {mix.JA, 5},
	"J1N":  {mix.J1, 0},
	"J1Z":  {mix.J1, 1},
	"J1P":  {mix.J1, 2},
	"J1NN": {mix.J1, 3},
	"J1NZ": {mix.J1, 4},
	"J1NP": {mix.J1, 5},
	"J2N":  {mix.J2, 0},
	"J2Z":  {mix.J2, 1},
	"J2P":  {mix.J2, 2},
	"J2NN": {mix.J2, 3},
	"J2NZ": {mix.J2, 4},
	"J2NP": {mix.J2, 5},
	"J3N":  {mix.J3, 0},
	"J3Z":  {mix.J3, 1},
	"J3P":  {mix.J3, 2},
	"J3NN": {mix.J3, 3},
	"J3NZ": {mix.J3, 4},
	"J3NP": {mix.J3, 5},
	"J4N":  {mix.J4, 0},
	"J4Z":  {mix.J4, 1},
	"J4P":  {mix.J4, 2},
	"J4NN": {mix.J4, 3},
	"J4NZ": {mix.J4, 4},
	"J4NP": {mix.J4, 5},
	"J5N":  {mix.J5, 0},
	"J5Z":  {mix.J5, 1},
	"J5P":  {mix.J5, 2},
	"J5NN": {mix.J5, 3},
	"J5NZ": {mix.J5, 4},
	"J5NP": {mix.J5, 5},
	"J6N":  {mix.J6, 0},
	"J6Z":  {mix.J6, 1},
	"J6P":  {mix.J6, 2},
	"J6NN": {mix.J6, 3},
	"J6NZ": {mix.J6, 4},
	"J6NP": {mix.J6, 5},
	"JXN":  {mix.JX, 0},
	"JXZ":  {mix.JX, 1},
	"JXP":  {mix.JX, 2},
	"JXNN": {mix.JX, 3},
	"JXNZ": {mix.JX, 4},
	"JXNP": {mix.JX, 5},
	"INCA": {mix.INCA, 0},
	"DECA": {mix.INCA, 1},
	"ENTA": {mix.INCA, 2},
	"ENNA": {mix.INCA, 3},
	"INC1": {mix.INC1, 0},
	"DEC1": {mix.INC1, 1},
	"ENT1": {mix.INC1, 2},
	"ENN1": {mix.INC1, 3},
	"INC2": {mix.INC2, 0},
	"DEC2": {mix.INC2, 1},
	"ENT2": {mix.INC2, 2},
	"ENN2": {mix.INC2, 3},
	"INC3": {mix.INC3, 0},
	"DEC3": {mix.INC3, 1},
	"ENT3": {mix.INC3, 2},
	"ENN3": {mix.INC3, 3},
	"INC4": {mix.INC4, 0},
	"DEC4": {mix.INC4, 1},
	"ENT4": {mix.INC4, 2},
	"ENN4": {mix.INC4, 3},
	"INC5": {mix.INC5, 0},
	"DEC5": {mix.INC5, 1},
	"ENT5": {mix.INC5, 2},
	"ENN5": {mix.INC5, 3},
	"INC6": {mix.INC6, 0},
	"DEC6": {mix.INC6, 1},
	"ENT6": {mix.INC6, 2},
	"ENN6": {mix.INC6, 3},
	"INCX": {mix.INCX, 0},
	"DECX": {mix.INCX, 1},
	"ENTX": {mix.INCX, 2},
	"ENNX": {mix.INCX, 3},
	"CMPA": {mix.CMPA, 5},
	"FCMP": {mix.CMPA, 6},
	"CMP1": {mix.CMP1, 5},
	"CMP2": {mix.CMP2, 5},
	"CMP3": {mix.CMP3, 5},
	"CMP4": {mix.CMP4, 5},
	"CMP5": {mix.CMP5, 5},
	"CMP6": {mix.CMP6, 5},
	"CMPX": {mix.CMPX, 5},
}
