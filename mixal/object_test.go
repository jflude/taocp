package mixal

import (
	"bytes"
	"testing"

	"github.com/jflude/gnuth/mix"
)

func TestFindAddress(t *testing.T) {
	w := egObject.findWord(1001)
	if w == nil {
		t.Error("got: nil, want: &1001")
	} else if w.Int() != -6882509 {
		t.Errorf("got: %d, want: -6882509", w.Int())
	}
	w = egObject.findWord(3002)
	if w == nil {
		t.Error("got: nil, want: &3002")
	} else if w.Int() != 133 {
		t.Errorf("got: %d, want: 133", w.Int())
	}
	w = egObject.findWord(2000)
	if w != nil {
		t.Errorf("got: %d, want: nil", w.Int())
	}
}

func TestWriteCards(t *testing.T) {
	var buf bytes.Buffer
	egObject.writeCards(&buf)
	if bytes.Compare(buf.Bytes(), okObject) != 0 {
		t.Errorf("got: %q\nwant: %q", buf.Bytes(), okObject)
	}
}

var egObject = object{
	start: 3000,
	orig:  []int{3000, 1000},
	seg: [][]mix.Word{
		{
			mix.NewWord(-1187),
			mix.NewWord(1000<<18 | 02245),
			mix.NewWord(133),
		},
		{
			mix.NewWord(135582544),
			mix.NewWord(-6882509),
			mix.NewWord(-67108864),
		},
	},
}

var okObject = []byte(" O O6 Y O6    I   B= D O4 Z IQ Z I3 Z EN    E   EU 0BB= H IU   EJ  CA. ACB=   EU 1A-H V A=  CEU 0AEH 1AEN    E  CLU  ABG H IH A A= J B. A  9                    ABCDE33000000000118P02621451890000000133                                        ABCDE310000135582544000688250R006710886M                                        TRANS03000                                                                      ")
