package mix

import (
	"errors"
	"io/ioutil"
	"strings"
	"testing"
)

func TestCycle(t *testing.T) {
	c, tmpDir := newSandbox(t, "")
	defer closeSandbox(t, c, tmpDir)
	c.trace = testing.Verbose()

	// LDA
	copy(c.Contents[mBase+0:], egCycle1)
	c.Contents[mBase+2000] = NewWord(-(80<<18 | 030504))
	c.next = 0
	for i, v := range okCycle1 {
		if err := c.Cycle(); err != nil {
			t.Errorf("#%d: error: %v", i+1, err)
			c.next++
			continue
		}
		if c.Reg[A] != v {
			t.Errorf("#%d: got: %#v, want: %#v",
				i+1, c.Reg[A], v)
		}
	}

	// LD[1-6]N
	copy(c.Contents[mBase+0:], egCycle2)
	c.Contents[mBase+2000] = NewWord(-01234)
	c.next = 0
	for i, v := range okCycle2 {
		if err := c.Cycle(); err != nil {
			t.Errorf("#%d: error: %v", i+1, err)
			c.next++
			continue
		}
		if c.Reg[I1+i] != v {
			t.Errorf("#%d: got: %#v, want: %#v",
				i+1, c.Reg[I1+i], v)
		}
	}

	// STA
	copy(c.Contents[mBase+0:], egCycle3)
	c.Reg[A] = NewWord(0607101100)
	c.next = 0
	for i, v := range okCycle3 {
		c.Contents[mBase+2000] = NewWord(-0102030405)
		if err := c.Cycle(); err != nil {
			t.Errorf("#%d: error: %v", i+1, err)
			c.next++
			continue
		}
		if c.Contents[mBase+2000] != v {
			t.Errorf("#%d: got: %#v, want: %#v",
				i+1, c.Contents[mBase+2000], v)
		}
	}

	// ST1
	copy(c.Contents[mBase+0:], egCycle4)
	c.Reg[I1] = NewWord(01100)
	c.next = 0
	for i, v := range okCycle4 {
		c.Contents[mBase+2000] = NewWord(-0102030405)
		if err := c.Cycle(); err != nil {
			t.Errorf("#%d: error: %v", i+1, err)
			c.next++
			continue
		}
		if c.Contents[mBase+2000] != v {
			t.Errorf("#%d: got: %#v, want: %#v",
				i+1, c.Contents[mBase+2000], v)
		}
	}

	// ADD, SUB, MUL, DIV
	for i, op := range egCycle5 {
		c.Contents[mBase+0] = op[0]
		c.Reg[A] = op[1]
		c.Reg[X] = op[2]
		c.Contents[mBase+1000] = op[3]
		c.next = 0
		if err := c.Cycle(); err != nil {
			t.Errorf("#%d: error: %v", i+1, err)
			continue
		}
		if c.Reg[A] != op[4] {
			t.Errorf("#%d: got: A = %#v, want: A = %#v",
				i+1, c.Reg[A], op[4])
		}
		if c.Reg[X] != op[5] {
			t.Errorf("#%d: got: X = %#v, want: X = %#v",
				i+1, c.Reg[X], op[5])
		}
	}

	// SLA, SRA, SRAX, SLC, SRC
	c.Reg[A] = NewWord(0102030405)
	c.Reg[X] = NewWord(-0607101112)
	c.next = 0
	for i, op := range egCycle6 {
		c.Contents[mBase+i] = op[0]
		if err := c.Cycle(); err != nil {
			t.Errorf("#%d: error: %v", i+1, err)
			c.next++
			continue
		}
		if c.Reg[A] != op[1] {
			t.Errorf("#%d: got: A = %#v, want: A = %#v",
				i+1, c.Reg[A], op[1])
			c.Reg[A] = op[1]
		}
		if c.Reg[X] != op[2] {
			t.Errorf("#%d: got: X = %#v, want: X = %#v",
				i+1, c.Reg[X], op[2])
			c.Reg[X] = op[2]
		}
	}

	// NUM, INCA1, CHAR
	c.Reg[A] = NewWord(-0374047)
	c.Reg[X] = NewWord(04571573636)
	c.next = 0
	for i, op := range egCycle7 {
		c.Contents[mBase+i] = op[0]
		if err := c.Cycle(); err != nil {
			t.Errorf("#%d: error: %v", i+1, err)
			c.next++
			continue
		}
		if c.Reg[A] != op[1] {
			t.Errorf("#%d: got: A = %#v, want: A = %#v",
				i+1, c.Reg[A], op[1])
			c.Reg[A] = op[1]
		}
		if c.Reg[X] != op[2] {
			t.Errorf("#%d: got: X = %#v, want: X = %#v",
				i+1, c.Reg[X], op[2])
			c.Reg[X] = op[2]
		}
	}

	// NUM (overflow, etc)
	c.Contents[mBase+0] = NUM
	for i, op := range egCycle8 {
		c.Reg[A] = op[0]
		c.Reg[X] = op[1]
		c.Overflow = false
		c.next = 0
		if err := c.Cycle(); err != nil {
			t.Errorf("#%d: error: %v", i+1, err)
			continue
		}
		if c.Overflow {
			c.Reg[A] = c.Reg[A].Negate()
		}
		if c.Reg[A] != op[2] {
			t.Errorf("#%d: got: A = %#v, want: A = %#v",
				i+1, c.Reg[A], op[2])
		}
	}

	// Program M, Section 1.3.2
	copy(c.Contents[mBase+3000:], egCycle9)
	c.Contents[mBase+0] = NewWord(3000<<18 | 39) // JMP 3000
	c.Contents[mBase+1] = NewWord(0205)          // HLT
	c.Contents[mBase+1000] = NewWord(1)
	c.Contents[mBase+1001] = NewWord(7)
	c.Contents[mBase+1002] = NewWord(2)
	c.Contents[mBase+1003] = NewWord(7)
	c.Contents[mBase+1004] = NewWord(6)
	c.Contents[mBase+1005] = NewWord(3)
	c.Contents[mBase+1006] = NewWord(1)
	c.Contents[mBase+1007] = NewWord(4)
	c.Contents[mBase+1008] = NewWord(5)
	c.Contents[mBase+1009] = NewWord(2)
	c.Reg[I1] = 10
	c.next = 0
	if err := c.resume(); !errors.Is(err, ErrHalted) {
		t.Error("error:", err)
	}
	if c.Reg[A].Int() != 7 {
		t.Errorf("got: %#o (%v), want: 7", c.Reg[A], c.Reg[A])
	}

	// Program P, Section 1.3.2
	for i := 0; i < len(c.Contents); i++ {
		c.Contents[i] = 0
	}
	copy(c.Contents[mBase+3000:], egCycle10)
	copy(c.Contents[mBase+0:], egCycle10a)
	copy(c.Contents[mBase+1995:], egCycle10b)
	copy(c.Contents[mBase+2024:], egCycle10c)
	copy(c.Contents[mBase+2049:], egCycle10d)
	c.next = 3000
	if err := c.resume(); !errors.Is(err, ErrHalted) {
		t.Error("error:", err)
	}
	b, err := ioutil.ReadFile("printer.mix")
	if err != nil {
		t.Fatal("error:", err)
	}
	if strings.Compare(string(b), okCycle10) != 0 {
		t.Error("got: incorrect printer output")
	}
}

func BenchmarkProgramM(b *testing.B) {
	c := NewComputer()
	copy(c.Contents[mBase+3000:], egCycle9)
	c.Contents[mBase+0] = NewWord(3000<<18 | 39) // JMP 3000
	c.Contents[mBase+1] = NewWord(0205)          // HLT
	c.Contents[mBase+1000] = NewWord(1)
	c.Contents[mBase+1001] = NewWord(7)
	c.Contents[mBase+1002] = NewWord(2)
	c.Contents[mBase+1003] = NewWord(7)
	c.Contents[mBase+1004] = NewWord(6)
	c.Contents[mBase+1005] = NewWord(3)
	c.Contents[mBase+1006] = NewWord(1)
	c.Contents[mBase+1007] = NewWord(4)
	c.Contents[mBase+1008] = NewWord(5)
	c.Contents[mBase+1009] = NewWord(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Reg[I1] = 10
		c.next = 0
		if err := c.resume(); !errors.Is(err, ErrHalted) {
			b.Fatal("error:", err)
		}
	}
}

func Benchmark1000Cycles(b *testing.B) {
	c := NewComputer()
	for i := 0; i < 998; i += 2 {
		c.Contents[mBase+i] = NewWord(0501)   // ADD 0
		c.Contents[mBase+i+1] = NewWord(0502) // SUB 0
	}
	c.Contents[mBase+999] = NewWord(0205) // HLT
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.next = 0
		if err := c.resume(); !errors.Is(err, ErrHalted) {
			b.Fatal("error:", err)
		}
	}
}

var (
	egCycle1 = []Word{
		NewWord(2000<<18 | 0510),  // LDA  2000
		NewWord(2000<<18 | 01510), // LDA  2000(1:5)
		NewWord(2000<<18 | 03510), // LDA  2000(3:5)
		NewWord(2000<<18 | 0310),  // LDA  2000(0:3)
		NewWord(2000<<18 | 04410), // LDA  2000(4:4)
		NewWord(2000<<18 | 010),   // LDA  2000(0:0)
	}
	okCycle1 = []Word{
		NewWord(-(80<<18 | 030504)),
		NewWord(80<<18 | 030504),
		NewWord(030504),
		NewWord(-(80<<6 | 3)),
		NewWord(5),
		NewWord(0).Negate(),
	}

	egCycle2 = []Word{
		NewWord(2000<<18 | 0521),  // LD1N 2000
		NewWord(2000<<18 | 01522), // LD2N 2000(1:5)
		NewWord(2000<<18 | 03523), // LD3N 2000(3:5)
		NewWord(2000<<18 | 0324),  // LD4N 2000(0:3)
		NewWord(2000<<18 | 04425), // LD5N 2000(4:4)
		NewWord(2000<<18 | 026),   // LD6N 2000(0:0)
	}
	okCycle2 = []Word{
		NewWord(01234),
		NewWord(-01234),
		NewWord(-01234),
		0,
		NewWord(-012),
		0,
	}

	egCycle3 = []Word{
		NewWord(2000<<18 | 0530),  // STA  2000
		NewWord(2000<<18 | 01530), // STA  2000(1:5)
		NewWord(2000<<18 | 05530), // STA  2000(5:5)
		NewWord(2000<<18 | 02230), // STA  2000(2:2)
		NewWord(2000<<18 | 02330), // STA  2000(2:3)
		NewWord(2000<<18 | 0130),  // STA  2000(0:1)
	}
	okCycle3 = []Word{
		NewWord(0607101100),
		NewWord(-0607101100),
		NewWord(-0102030400),
		NewWord(-0100030405),
		NewWord(-0111000405),
		NewWord(02030405),
	}

	egCycle4 = []Word{
		NewWord(2000<<18 | 0531),  // ST1  2000
		NewWord(2000<<18 | 01531), // ST1  2000(1:5)
		NewWord(2000<<18 | 05531), // ST1  2000(5:5)
		NewWord(2000<<18 | 02231), // ST1  2000(2:2)
		NewWord(2000<<18 | 02331), // ST1  2000(2:3)
		NewWord(2000<<18 | 0131),  // ST1  2000(0:1)
	}
	okCycle4 = []Word{
		NewWord(01100),
		NewWord(-01100),
		NewWord(-0102030400),
		NewWord(-0100030405),
		NewWord(-0111000405),
		NewWord(02030405),
	}

	egCycle5 = [][]Word{
		[]Word{ // #1
			NewWord(01750000501),       // ADD  1000
			NewWord(1234<<18 | 010226), // A (before)
			0,                          // X (before)
			NewWord(100<<18 | 050062),  // CONTENTS[1000]
			NewWord(1334<<18 | 060310), // A (after)
			0,                          // X (after)
		},
		[]Word{ // #2
			NewWord(01750000502), // SUB  1000
			NewWord(-(1234<<18 | 9)),
			0,
			NewWord(-(2000<<18 | (150 << 6))),
			NewWord(766<<18 | 149<<6 | 067),
			0,
		},
		[]Word{ // #3
			NewWord(01750001103), // MUL  1000(1:1)
			NewWord(-112),
			0,
			NewWord(0200000000),
			NewWord(0).Negate(),
			NewWord(-224),
		},
		[]Word{ // #4
			NewWord(01750000503), // MUL  1000
			NewWord(-(50<<24 | 112<<6 | 4)),
			0,
			NewWord(-0200000000),
			NewWord(100<<18 | 224),
			NewWord(8 << 24),
		},
		[]Word{ // #5
			NewWord(01750000504), // DIV  1000
			0,
			NewWord(17),
			NewWord(3),
			NewWord(5),
			NewWord(2),
		},
		[]Word{ // #6
			NewWord(01750000504), // DIV  1000
			NewWord(0).Negate(),
			NewWord(1235<<18 | 0301),
			NewWord(-0200),
			NewWord(617<<12 | 04001),
			NewWord(-0101),
		},
	}

	egCycle6 = [][3]Word{
		{ // #1
			NewWord(01000306),    // SRAX 1
			NewWord(01020304),    // A
			NewWord(-0506071011), // X
		},
		{ // #2
			NewWord(02000006), // SLA  2
			NewWord(0203040000),
			NewWord(-0506071011),
		},
		{ // #3
			NewWord(04000506), // SRC  4
			NewWord(0607101102),
			NewWord(-0304000005),
		},
		{ // #4
			NewWord(02000106), // SRA  2
			NewWord(060710),
			NewWord(-0304000005),
		},
		{ // #5
			NewWord(0765000406), // SLC  501
			NewWord(06071003),
			NewWord(-0400000500),
		},
	}

	egCycle7 = [][3]Word{
		{ // #1
			NewWord(05), // NUM
			NewWord(-12977700),
			NewWord(04571573636),
		},
		{ // #2
			NewWord(01000060), // INCA 1
			NewWord(-12977699),
			NewWord(04571573636),
		},
		{ // #3
			NewWord(0105), // CHAR
			NewWord(-03636374047),
			NewWord(04545444747),
		},
	}

	egCycle8 = [][3]Word{
		{ // #1
			NewWord(-03736363636), // A (before)
			NewWord(-03636363637), // X (before)
			NewWord(-1000000001),  // A (after, negated if overflow)
		},
		{ // #2
			NewWord(-04747474747),
			NewWord(04747474747),
			NewWord(02402761777),
		},
	}

	//                                         * FIND THE MAXIMUM
	//                                         X        EQU  1000
	egCycle9 = []Word{ //                               ORIG 3000
		NewWord(3009<<18 | 0240),       // MAXIMUM  STJ  EXIT
		NewWord(010263),                // INIT     ENT3 0,1
		NewWord(3005<<18 | 39),         //          JMP  CHANGEM
		NewWord(1000<<18 | 030570),     // LOOP     CMPA X,3
		NewWord(3007<<18 | 0700 | 39),  //          JLE  *+3
		NewWord(030200 | 50),           // CHANGEM  ENT2 0,3
		NewWord(1000<<18 | 030500 | 8), //          LDA  X,3
		NewWord(01000100 | 51),         //          DEC3 1
		NewWord(3003<<18 | 0200 | 43),  //          J3P  LOOP
		NewWord(3009<<18 | 39),         // EXIT     JMP  *
	}

	//                                    * EXAMPLE: TABLE OF PRIMES
	//                                    L         EQU  500
	//                                    PRINTER   EQU  18
	//                                    PRIME     EQU  -1
	//                                    BUF0      EQU  2000
	//                                    BUF1      EQU  BUF0+25
	egCycle10 = []Word{ //                          ORIG 3000
		NewWord(02243),            // START     IOC  0(PRINTER)
		NewWord(2050<<18 | 0511),  //           LD1  =1-L=
		NewWord(2051<<18 | 0512),  //           LD2  =3=
		NewWord(01000061),         // 2H        INC1 1
		NewWord(499<<18 | 010532), //           ST2  PRIME+L,1
		NewWord(3016<<18 | 0151),  //           J1Z  2F
		NewWord(02000062),         // 4H        INC2 2
		NewWord(02000263),         //           ENT3 2
		NewWord(0260),             // 6H        ENTA 0
		NewWord(020267),           //           ENTX 0,2
		NewWord(-01030504),        //           DIV  PRIME,3
		NewWord(3006<<18 | 0157),  //           JXZ  4B
		NewWord(-01030570),        //           CMPA PRIME,3
		NewWord(01000063),         //           INC3 1
		NewWord(3008<<18 | 0647),  //           JG   6B
		NewWord(3003<<18 | 047),   //           JMP  2B
		NewWord(1995<<18 | 02245), // 2H        OUT  TITLE(PRINTER)
		NewWord(2035<<18 | 0264),  //           ENT4 BUF1+10
		NewWord(-062000265),       //           ENT5 -50
		NewWord(501<<18 | 065),    // 2H        INC5 L+1
		NewWord(-01050510),        // 4H        LDA  PRIME,5
		NewWord(0105),             //           CHAR
		NewWord(041437),           //           STX  0,4(1:4)
		NewWord(01000164),         //           DEC4 1
		NewWord(062000165),        //           DEC5 50
		NewWord(3020<<18 | 0255),  //           J5P  4B
		NewWord(042245),           //           OUT  0,4(PRINTER)
		NewWord(030040514),        //           LD4  24,4
		NewWord(3019<<18 | 055),   //           J5N  2B
		NewWord(0205),             //           HLT
	}
	//                                    * TABLES AND BUFFERS
	egCycle10a = []Word{ //                         ORIG PRIME+1 (=0)
		NewWord(2), //                          CON  2
	}
	egCycle10b = []Word{ //                         ORIG BUF0-5 (=1995)
		NewWord(0611232627),  //      TITLE     ALF  FIRST
		NewWord(06113105),    //                ALF   FIVE
		NewWord(010301704),   //                ALF   HUND
		NewWord(02305040021), //                ALF  RED P
		NewWord(02311160526), //                ALF  RIMES
	}
	egCycle10c = []Word{ //                         ORIG BUF0+24 (=2024)
		NewWord(2035), //                       CON  BUF1+10
	}
	egCycle10d = []Word{ //                         ORIG BUF1+24 (=2049)
		NewWord(2010), //                       CON  BUF0+10
		NewWord(-499), //                       CON  1-L
		NewWord(3),    //                       CON  3
	}

	okCycle10 = "\014" + `FIRST FIVE HUNDRED PRIMES
     0002 0233 0547 0877 1229 1597 1993 2371 2749 3187
     0003 0239 0557 0881 1231 1601 1997 2377 2753 3191
     0005 0241 0563 0883 1237 1607 1999 2381 2767 3203
     0007 0251 0569 0887 1249 1609 2003 2383 2777 3209
     0011 0257 0571 0907 1259 1613 2011 2389 2789 3217
     0013 0263 0577 0911 1277 1619 2017 2393 2791 3221
     0017 0269 0587 0919 1279 1621 2027 2399 2797 3229
     0019 0271 0593 0929 1283 1627 2029 2411 2801 3251
     0023 0277 0599 0937 1289 1637 2039 2417 2803 3253
     0029 0281 0601 0941 1291 1657 2053 2423 2819 3257
     0031 0283 0607 0947 1297 1663 2063 2437 2833 3259
     0037 0293 0613 0953 1301 1667 2069 2441 2837 3271
     0041 0307 0617 0967 1303 1669 2081 2447 2843 3299
     0043 0311 0619 0971 1307 1693 2083 2459 2851 3301
     0047 0313 0631 0977 1319 1697 2087 2467 2857 3307
     0053 0317 0641 0983 1321 1699 2089 2473 2861 3313
     0059 0331 0643 0991 1327 1709 2099 2477 2879 3319
     0061 0337 0647 0997 1361 1721 2111 2503 2887 3323
     0067 0347 0653 1009 1367 1723 2113 2521 2897 3329
     0071 0349 0659 1013 1373 1733 2129 2531 2903 3331
     0073 0353 0661 1019 1381 1741 2131 2539 2909 3343
     0079 0359 0673 1021 1399 1747 2137 2543 2917 3347
     0083 0367 0677 1031 1409 1753 2141 2549 2927 3359
     0089 0373 0683 1033 1423 1759 2143 2551 2939 3361
     0097 0379 0691 1039 1427 1777 2153 2557 2953 3371
     0101 0383 0701 1049 1429 1783 2161 2579 2957 3373
     0103 0389 0709 1051 1433 1787 2179 2591 2963 3389
     0107 0397 0719 1061 1439 1789 2203 2593 2969 3391
     0109 0401 0727 1063 1447 1801 2207 2609 2971 3407
     0113 0409 0733 1069 1451 1811 2213 2617 2999 3413
     0127 0419 0739 1087 1453 1823 2221 2621 3001 3433
     0131 0421 0743 1091 1459 1831 2237 2633 3011 3449
     0137 0431 0751 1093 1471 1847 2239 2647 3019 3457
     0139 0433 0757 1097 1481 1861 2243 2657 3023 3461
     0149 0439 0761 1103 1483 1867 2251 2659 3037 3463
     0151 0443 0769 1109 1487 1871 2267 2663 3041 3467
     0157 0449 0773 1117 1489 1873 2269 2671 3049 3469
     0163 0457 0787 1123 1493 1877 2273 2677 3061 3491
     0167 0461 0797 1129 1499 1879 2281 2683 3067 3499
     0173 0463 0809 1151 1511 1889 2287 2687 3079 3511
     0179 0467 0811 1153 1523 1901 2293 2689 3083 3517
     0181 0479 0821 1163 1531 1907 2297 2693 3089 3527
     0191 0487 0823 1171 1543 1913 2309 2699 3109 3529
     0193 0491 0827 1181 1549 1931 2311 2707 3119 3533
     0197 0499 0829 1187 1553 1933 2333 2711 3121 3539
     0199 0503 0839 1193 1559 1949 2339 2713 3137 3541
     0211 0509 0853 1201 1567 1951 2341 2719 3163 3547
     0223 0521 0857 1213 1571 1973 2347 2729 3167 3557
     0227 0523 0859 1217 1579 1979 2351 2731 3169 3559
     0229 0541 0863 1223 1583 1987 2357 2741 3181 3571
`
)
