// Copyright (c) 2020 Justin Flude. All rights reserved.
// Use of this source code is governed by the COPYING.md file.
package mix

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestInterrupt(t *testing.T) {
	c, tmpDir := newSandbox(t, "")
	defer closeSandbox(t, c, tmpDir)
	if testing.Verbose() {
		c.Tracer = os.Stdout
	}

	copy(c.Contents[mBase-1000:], egInterrupt1)
	copy(c.Contents[mBase-11:], egInterrupt2)
	copy(c.Contents[mBase-38:], egInterrupt3)
	copy(c.Contents[mBase+1000:], egInterrupt4)
	c.next = 1000
	c.Interrupts = true
	if err := c.resume(); !errors.Is(err, ErrHalted) {
		t.Error("error:", err)
	}
	if c.Elapsed != 12100021 {
		t.Errorf("got: %du elapsed, want: 12100021u", c.Elapsed)
	}
	b, err := os.ReadFile("printer.mix")
	if err != nil {
		t.Fatal("error:", err)
	}
	if strings.Compare(string(b), okInterrupt) != 0 {
		t.Error("got: incorrect printer output")
	}
}

var egInterrupt1 = []Word{ //     * EXAMPLE: INTERRUPT CAPABILITY
	//                        * SUPERVISOR ROUTINE
	//                        QUANTUM    EQU  1000
	//                        CLOCK      EQU  -10
	//                        PRINTER    EQU  18
	//                        STATUS     EQU  -1
	//                        NEXTOP     EQU  0:2
	//                        BUF        EQU  -2000
	//                                   ORIG -1000
	NewWord(01750000260),  // SUPER      ENTA QUANTUM
	NewWord(-012000530),   //            STA  CLOCK
	NewWord(-01741002242), //            JBUS RESUME(PRINTER)
	NewWord(-01000210),    //            LDA  STATUS(NEXTOP)
	NewWord(0105),         //            CHAR
	NewWord(-03720000537), //            STX  BUF
	NewWord(-03720002245), //            OUT  BUF(PRINTER)
	NewWord(01105),        // RESUME     INT
}

var egInterrupt2 = []Word{ //     * INTERRUPT VECTORS
	//                        CLOCKVEC   EQU  -11
	//                        PRINTVEC   EQU  -20-PRINTER
	//                                   ORIG CLOCKVEC
	NewWord(-01750000047), //            JMP  SUPER
	NewWord(01750),        //            CON  QUANTUM
}

var egInterrupt3 = []Word{ //                ORIG PRINTVEC
	NewWord(01105), //                   INT
}

var egInterrupt4 = []Word{ //     * WORKER ROUTINE
	//                                   ORIG 1000
	NewWord(01766000510), //  START      LDA  =1000000=
	NewWord(0500),        //  LOOP       NOP
	NewWord(0500),        //             NOP
	NewWord(0500),        //             NOP
	NewWord(0500),        //             NOP
	NewWord(0500),        //             NOP
	NewWord(0500),        //             NOP
	NewWord(0500),        //             NOP
	NewWord(0500),        //             NOP
	NewWord(0500),        //             NOP
	NewWord(0500),        //             NOP
	NewWord(01000160),    //             DECA 1
	NewWord(01751000450), //             JANZ LOOP
	NewWord(0205),        //             HLT
	NewWord(03641100),    //             END  START
}

var okInterrupt = `01003
01007
01011
01003
01007
01011
01003
01007
01011
01003
01007
01011
`
