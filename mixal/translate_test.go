// Copyright (c) 2020 Justin Flude. All rights reserved.
// Use of this source code is governed by the COPYING.md file.
package mixal

import (
	"strings"
	"testing"
)

func TestTranslate(t *testing.T) {
	r := strings.NewReader(egTranslate)
	var a asmb
	if err := a.translate(r, func(a *asmb, loc, op, address string) {
		if loc != okTranslate[a.count].loc ||
			op != okTranslate[a.count].op ||
			address != okTranslate[a.count].address {
			t.Errorf("%d: got: %q, %q, %q, want: %v",
				a.count, loc, op, address,
				okTranslate[a.count])
		}
	}); err != nil {
		t.Error("error:", err)
	}
}

var egTranslate = `* EXAMPLE: TABLE OF PRIMES
L          EQU  500
PRINTER    EQU  18
PRIME      EQU  -1
BUF0       EQU  2000
BUF1       EQU  BUF0+25
           ORIG 3000
START      IOC  0(PRINTER)
           LD1  =1-L=
           LD2  =3=
2H         INC1 1
           ST2  PRIME+L,1
           J1Z  2F
4H         INC2 2
           ENT3 2
6H         ENTA 0
           ENTX 0,2
           DIV  PRIME,3
           JXZ  4B
           CMPA PRIME,3
           INC3 1
           JG   6B
           JMP  2B
2H         OUT  TITLE(PRINTER)
           ENT4 BUF1+10
           ENT5 -50
2H         INC5 L+1
4H         LDA  PRIME,5
           CHAR
           STX  0,4(1:4)
           DEC4 1
           DEC5 50
           J5P  4B
           OUT  0,4(PRINTER)
           LD4  24,4
           J5N  2B
           HLT
* TABLES AND BUFFERS
           ORIG PRIME+1
           CON  2
           ORIG BUF0-5
TITLE      ALF  FIRST
           ALF   FIVE
           ALF   HUND
           ALF  RED P
           ALF  RIMES
           ORIG BUF0+24
           CON  BUF1+10
           ORIG BUF1+24
           CON  BUF0+10
           END  START
`

var okTranslate = []struct {
	loc, op, address string
}{
	2:  {"L", "EQU", "500"},
	3:  {"PRINTER", "EQU", "18"},
	4:  {"PRIME", "EQU", "-1"},
	5:  {"BUF0", "EQU", "2000"},
	6:  {"BUF1", "EQU", "BUF0+25"},
	7:  {"", "ORIG", "3000"},
	8:  {"START", "IOC", "0(PRINTER)"},
	9:  {"", "LD1", "=1-L="},
	10: {"", "LD2", "=3="},
	11: {"2H", "INC1", "1"},
	12: {"", "ST2", "PRIME+L,1"},
	13: {"", "J1Z", "2F"},
	14: {"4H", "INC2", "2"},
	15: {"", "ENT3", "2"},
	16: {"6H", "ENTA", "0"},
	17: {"", "ENTX", "0,2"},
	18: {"", "DIV", "PRIME,3"},
	19: {"", "JXZ", "4B"},
	20: {"", "CMPA", "PRIME,3"},
	21: {"", "INC3", "1"},
	22: {"", "JG", "6B"},
	23: {"", "JMP", "2B"},
	24: {"2H", "OUT", "TITLE(PRINTER)"},
	25: {"", "ENT4", "BUF1+10"},
	26: {"", "ENT5", "-50"},
	27: {"2H", "INC5", "L+1"},
	28: {"4H", "LDA", "PRIME,5"},
	29: {"", "CHAR", ""},
	30: {"", "STX", "0,4(1:4)"},
	31: {"", "DEC4", "1"},
	32: {"", "DEC5", "50"},
	33: {"", "J5P", "4B"},
	34: {"", "OUT", "0,4(PRINTER)"},
	35: {"", "LD4", "24,4"},
	36: {"", "J5N", "2B"},
	37: {"", "HLT", ""},
	39: {"", "ORIG", "PRIME+1"},
	40: {"", "CON", "2"},
	41: {"", "ORIG", "BUF0-5"},
	42: {"TITLE", "ALF", "FIRST"},
	43: {"", "ALF", " FIVE"},
	44: {"", "ALF", " HUND"},
	45: {"", "ALF", "RED P"},
	46: {"", "ALF", "RIMES"},
	47: {"", "ORIG", "BUF0+24"},
	48: {"", "CON", "BUF1+10"},
	49: {"", "ORIG", "BUF1+24"},
	50: {"", "CON", "BUF0+10"},
	51: {"", "END", "START"},
}
