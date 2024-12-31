package mixal

import (
	"reflect"
	"strings"
	"testing"

	"github.com/jflude/taocp/mix"
)

func TestEvaluating(t *testing.T) {
	checkParsing(t, egParsing1, &okParsing1)
}

func TestAssembling(t *testing.T) {
	checkParsing(t, egTranslate, &okParsing2)
}

func checkParsing(t *testing.T, data string, expected *object) {
	r := strings.NewReader(data)
	lastEmit := 0
	var a asmb
	if err := a.translate(r, func(a *asmb, loc, op, address string) {
		t.Logf("input: %d: %s %s %s", a.count, loc, op, address)
		parseLine(a, loc, op, address)
		t.Log("  token:", a.tokens)
		if a.obj.orig == nil {
			return
		}
		orig := a.obj.orig[len(a.obj.seg)-1]
		seg := a.obj.seg[len(a.obj.seg)-1]
		if lastEmit > len(seg)-1 {
			lastEmit = 0
		}
		for ; lastEmit < len(seg); lastEmit++ {
			t.Logf("  emit: %d: %#v (%v)", orig+lastEmit,
				seg[lastEmit], seg[lastEmit])
		}
	}); err != nil {
		t.Fatal("error:", err)
	}
	if a.obj.start != expected.start {
		t.Errorf("object.start: got: %v, want: %v",
			a.obj.start, expected.start)
	}
	if !reflect.DeepEqual(a.obj.orig, expected.orig) {
		t.Errorf("object.orig: got: %v, want: %v",
			a.obj.orig, expected.orig)
	}
outer:
	for i := range expected.seg {
		if i >= len(a.obj.seg) {
			t.Errorf("seg[%d]: got: nil", i)
			break
		}
		for j := range expected.seg[i] {
			if j >= len(a.obj.seg[i]) {
				t.Errorf("seg[%d][%d]: got: nil", i, j)
				break outer
			}
			if a.obj.seg[i][j] != expected.seg[i][j] {
				t.Errorf("seg[%d][%d]: got: %#v, want: %#v",
					i, j, a.obj.seg[i][j],
					expected.seg[i][j])
				break outer
			}
		}
	}
}

var egParsing1 = `* TEST EXPRESSION EVALUATION
           ORIG 1000
START      HLT  *+3
           CON  9(1:1),18(2:2),27(3:3),36(4:4),45(5:5)
           CON  63(2:2),22/11
           CON  -3*8
           CON  10/-2
           CON  1//-3
           CON  -1+5*20/6+*
           END  START
`

var okParsing1 = object{
	start: 1000,
	orig:  []int{1000},
	seg: [][]mix.Word{
		[]mix.Word{
			mix.NewWord(01753000205),
			mix.NewWord(01122334455),
			mix.NewWord(2),
			mix.NewWord(-24),
			mix.NewWord(-5),
			mix.NewWord(-02525252525),
			mix.NewWord(1019),
		},
	},
}

var okParsing2 = object{
	start: 3000,
	orig:  []int{3000, 0, 1995, 2024, 2049},
	seg: [][]mix.Word{
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
