package mixal

import (
	"reflect"
	"strings"
	"testing"

	"github.com/jflude/gnuth/mix"
)

func TestParsing(t *testing.T) {
	r := strings.NewReader(egTranslate)
	var a asmb
	var lastWord *mix.Word
	if err := a.translate(r, func(a *asmb, loc, op, address string) {
		t.Logf("input: %d: %s %s %s", a.count, loc, op, address)
		parseLine(a, loc, op, address)
		t.Log("token:", a.tokens)
		seg := a.obj.seg[len(a.obj.seg)-1]
		if seg != nil && lastWord != &seg[len(seg)-1] {
			t.Logf(" emit: %d: %#o",
				a.obj.orig[len(a.obj.seg)-1]+len(seg)-1,
				seg[len(seg)-1])
			lastWord = &seg[len(seg)-1]
		}
		t.Log()
	}); err != nil {
		t.Fatal("error:", err)
	}
	if a.obj.start != okParsing.start {
		t.Errorf("object.start: got: %v, want: %v",
			a.obj.start, okParsing.start)
	}
	if !reflect.DeepEqual(a.obj.orig, okParsing.orig) {
		t.Errorf("object.orig: got: %v, want: %v",
			a.obj.orig, okParsing.orig)
	}
outer:
	for i := 0; i < len(okParsing.seg); i++ {
		if i >= len(a.obj.seg) {
			t.Errorf("seg[%d]: got: nil", i)
			break
		}
		for j := 0; j < len(okParsing.seg[i]); j++ {
			if j >= len(a.obj.seg[i]) {
				t.Errorf("seg[%d][%d]: got: nil", i, j)
				break outer
			}
			if a.obj.seg[i][j] != okParsing.seg[i][j] {
				t.Errorf("seg[%d][%d]: got: %#o, want: %#o",
					i, j, a.obj.seg[i][j],
					okParsing.seg[i][j])
				break outer
			}
		}
	}
}

var okParsing = object{
	start: 3000,
	orig:  []int{0, 3000, 0, 1995, 2024, 2049},
	seg: [][]mix.Word{
		nil,
		//                                        * EXAMPLE: TABLE OF PRIMES
		//                                        L         EQU   500
		//                                        PRINTER   EQU   18
		//                                        PRIME     EQU   -1
		//                                        BUF0      EQU   2000
		//                                        BUF1      EQU   BUF0+25
		[]mix.Word{ //                                      ORIG  3000
			mix.NewWord(02243),            // START     IOC   0(PRINTER)
			mix.NewWord(2050<<18 | 0511),  //           LD1   =1-L=
			mix.NewWord(2051<<18 | 0512),  //           LD2   =3=
			mix.NewWord(01000061),         // 2H        INC1  1
			mix.NewWord(499<<18 | 010532), //           ST2   PRIME+L,1
			mix.NewWord(3016<<18 | 0151),  //           J1Z   2F
			mix.NewWord(02000062),         // 4H        INC2  2
			mix.NewWord(02000263),         //           ENT3  2
			mix.NewWord(0260),             // 6H        ENTA  0
			mix.NewWord(020267),           //           ENTX  0,2
			mix.NewWord(-01030504),        //           DIV   PRIME,3
			mix.NewWord(3006<<18 | 0157),  //           JXZ   4B
			mix.NewWord(-01030570),        //           CMPA  PRIME,3
			mix.NewWord(01000063),         //           INC3  1
			mix.NewWord(3008<<18 | 0647),  //           JG    6B
			mix.NewWord(3003<<18 | 047),   //           JMP   2B
			mix.NewWord(1995<<18 | 02245), // 2H        OUT   TITLE(PRINTER)
			mix.NewWord(2035<<18 | 0264),  //           ENT4  BUF1+10
			mix.NewWord(-062000265),       //           ENT5  -50
			mix.NewWord(501<<18 | 065),    // 2H        INC5  L+1
			mix.NewWord(-01050510),        // 4H        LDA   PRIME,5
			mix.NewWord(0105),             //           CHAR
			mix.NewWord(041437),           //           STX   0,4(1:4)
			mix.NewWord(01000164),         //           DEC4  1
			mix.NewWord(062000165),        //           DEC5  50
			mix.NewWord(3020<<18 | 0255),  //           J5P   4B
			mix.NewWord(042245),           //           OUT   0,4(PRINTER)
			mix.NewWord(030040514),        //           LD4   24,4
			mix.NewWord(3019<<18 | 055),   //           J5N   2B
			mix.NewWord(0205),             //           HLT
		},
		//                                        * TABLES AND BUFFERS
		[]mix.Word{ //                                      ORIG  PRIME+1 (=0)
			mix.NewWord(2), //                          CON   2
		},
		[]mix.Word{ //                                      ORIG  BUF0-5 (=1995)
			mix.NewWord(0611232627),  //      TITLE     ALF   FIRST
			mix.NewWord(06113105),    //                ALF    FIVE
			mix.NewWord(010301704),   //                ALF    HUND
			mix.NewWord(02305040021), //                ALF   RED P
			mix.NewWord(02311160526), //                ALF   RIMES
		},
		[]mix.Word{ //                                      ORIG  BUF0+24 (=2024)
			mix.NewWord(2035), //                       CON   BUF1+10
		},
		[]mix.Word{ //                                      ORIG  BUF1+24 (=2049)
			mix.NewWord(2010), //                       CON   BUF0+10
			mix.NewWord(-499), //                       CON   1-L
			mix.NewWord(3),    //                       CON   3
		},
	},
}
