package mix

import "testing"

func TestCycle(t *testing.T) {
	c := NewComputer()

	// LDA
	copy(c.Contents[:], egCycle1)
	c.Contents[2000] = NewWord(-(80<<18 | 030504))
	c.next = 0
	for i, v := range okCycle1 {
		if err := c.Cycle(); err != nil {
			t.Errorf("#%d: got error: %v", i+1, err)
			c.next++
			break
		}
		if c.Reg[A] != v {
			t.Errorf("#%d: got %#v, want %#v",
				i+1, c.Reg[A], v)
		}
	}

	// LD[1-6]N
	copy(c.Contents[:], egCycle2)
	c.Contents[2000] = NewWord(-01234)
	c.next = 0
	for i, v := range okCycle2 {
		if err := c.Cycle(); err != nil {
			t.Errorf("#%d: got error: %v", i+1, err)
			c.next++
			continue
		}
		if c.Reg[I1+i] != v {
			t.Errorf("#%d: got %#v, want %#v",
				i+1, c.Reg[I1+i], v)
		}
	}

	// STA
	copy(c.Contents[:], egCycle3)
	c.Reg[A] = NewWord(0607101100)
	c.next = 0
	for i, v := range okCycle3 {
		c.Contents[2000] = NewWord(-0102030405)
		if err := c.Cycle(); err != nil {
			t.Errorf("#%d: got error: %v", i+1, err)
			c.next++
			continue
		}
		if c.Contents[2000] != v {
			t.Errorf("#%d: got %#v, want %#v",
				i+1, c.Contents[2000], v)
		}
	}

	// ST1
	copy(c.Contents[:], egCycle4)
	c.Reg[I1] = NewWord(01100)
	c.next = 0
	for i, v := range okCycle4 {
		c.Contents[2000] = NewWord(-0102030405)
		if err := c.Cycle(); err != nil {
			t.Errorf("#%d: got error: %v", i+1, err)
			c.next++
			continue
		}
		if c.Contents[2000] != v {
			t.Errorf("#%d: got %#v, want %#v",
				i+1, c.Contents[2000], v)
		}
	}

	// ADD, SUB, MUL, DIV
	for i, op := range egCycle5 {
		c.Contents[0] = op[0]
		c.Reg[A] = op[1]
		c.Reg[X] = op[2]
		c.Contents[1000] = op[3]
		c.next = 0
		if err := c.Cycle(); err != nil {
			t.Errorf("#%d: got error: %v", i+1, err)
			continue
		}
		if c.Reg[A] != op[4] {
			t.Errorf("#%d: got A = %#v, want A = %#v",
				i+1, c.Reg[A], op[4])
		}
		if c.Reg[X] != op[5] {
			t.Errorf("#%d: got X = %#v, want X = %#v",
				i+1, c.Reg[X], op[5])
		}
	}
}

var (
	egCycle1 = []Word{
		NewWord(2000<<18 | 0510),  // LDA 2000
		NewWord(2000<<18 | 01510), // LDA 2000(1:5)
		NewWord(2000<<18 | 03510), // LDA 2000(3:5)
		NewWord(2000<<18 | 0310),  // LDA 2000(0:3)
		NewWord(2000<<18 | 04410), // LDA 2000(4:4)
		NewWord(2000<<18 | 010),   // LDA 2000(0:0)
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
		NewWord(0),
		NewWord(-012),
		NewWord(0),
	}

	egCycle3 = []Word{
		NewWord(2000<<18 | 0530),  // STA 2000
		NewWord(2000<<18 | 01530), // STA 2000(1:5)
		NewWord(2000<<18 | 05530), // STA 2000(5:5)
		NewWord(2000<<18 | 02230), // STA 2000(2:2)
		NewWord(2000<<18 | 02330), // STA 2000(2:3)
		NewWord(2000<<18 | 0130),  // STA 2000(0:1)
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
		NewWord(2000<<18 | 0531),  // ST1 2000
		NewWord(2000<<18 | 01531), // ST1 2000(1:5)
		NewWord(2000<<18 | 05531), // ST1 2000(5:5)
		NewWord(2000<<18 | 02231), // ST1 2000(2:2)
		NewWord(2000<<18 | 02331), // ST1 2000(2:3)
		NewWord(2000<<18 | 0131),  // ST1 2000(0:1)
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
			NewWord(01750000501),       // ADD 1000
			NewWord(1234<<18 | 010226), // A
			NewWord(0),                 // X
			NewWord(100<<18 | 050062),  // CONTENTS[1000]
			NewWord(1334<<18 | 060310), // A (after)
			NewWord(0),                 // X (after)
		},
		[]Word{ // #2
			NewWord(01750000502), // SUB 1000
			NewWord(-(1234<<18 | 9)),
			NewWord(0),
			NewWord(-(2000<<18 | (150 << 6))),
			NewWord(766<<18 | 149<<6 | 067),
			NewWord(0),
		},
		[]Word{ // #3
			NewWord(01750001103), // MUL 1000(1:1)
			NewWord(-112),
			NewWord(0),
			NewWord(0200000000),
			NewWord(0).Negate(),
			NewWord(-224),
		},
		[]Word{ // #4
			NewWord(01750000503), // MUL 1000
			NewWord(-(50<<24 | 112<<6 | 4)),
			NewWord(0),
			NewWord(-0200000000),
			NewWord(100<<18 | 224),
			NewWord(8 << 24),
		},
		[]Word{ // #5
			NewWord(01750000504), // DIV 1000
			NewWord(0),
			NewWord(17),
			NewWord(3),
			NewWord(5),
			NewWord(2),
		},
		[]Word{ // #6
			NewWord(01750000504), // DIV 1000
			NewWord(0).Negate(),
			NewWord(1235<<18 | 0301),
			NewWord(-0200),
			NewWord(617<<12 | 04001),
			NewWord(-0101),
		},
	}
)
