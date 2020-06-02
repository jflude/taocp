package mixal

import (
	"bytes"
	"testing"

	"github.com/jflude/gnuth/mix"
)

func TestOutput(t *testing.T) {
	var buf bytes.Buffer
	egObject.outputCards(&buf, 3000)
	if bytes.Compare(buf.Bytes(), okObject) != 0 {
		t.Errorf("got: %q\nwant: %q", buf.Bytes(), okObject)
	}
}

var egObject = object{
	&chunk{
		orig: 3000,
		data: []mix.Word{
			mix.NewWord(1187),
			mix.NewWord(1000<<18 | 02245),
			mix.NewWord(133),
		},
	},
	&chunk{
		orig: 1000,
		data: []mix.Word{
			mix.NewWord(135582544),
			mix.NewWord(6882509),
			mix.NewWord(67108864),
		},
	},
}

var okObject = []byte(" O O6 Y O6    I   B= D O4 Z IQ Z I3 Z EN    E   EU 0BB= H IU   EJ  CA. ACB=   EU 1A-H V A=  CEU 0AEH 1AEN    E  CLU  ABG H IH A A= J B. A  9                    ABCDE33000000000118702621451890000000133                                        ABCDE31000013558254400068825090067108864                                        TRANS03000                                                                      ")
