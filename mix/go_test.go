package mix

import (
	"errors"
	"io/ioutil"
	"strings"
	"testing"
)

func TestGoButton(t *testing.T) {
	c, tmpDir := newSandbox(t, "")
	defer closeSandbox(t, c, tmpDir)
	c.trace = testing.Verbose()

	deck := append(append([]byte(nil), egLoader[0]...), egLoader[1]...)
	deck = append(append(deck, egHelloWorld...), egTransfer...)
	err := ioutil.WriteFile("reader.mix", deck, 0644)
	if err != nil {
		t.Fatal("error:", err)
	}
	if err = c.GoButton(16); !errors.Is(err, ErrHalted) {
		t.Error("error:", err)
	}
	b, err := ioutil.ReadFile("printer.mix")
	if err != nil {
		t.Fatal("error:", err)
	}
	if strings.Compare(string(b), okLoader) != 0 {
		t.Error("got: incorrect printer output")
	}
}

var egLoader = []string{
	//   5   10   15   20   25   30   35   40   45   50   55   60   65   70   75   80
	" O O6 Y O6    I   B= D O4 Z IQ Z I3 Z EN    E   EU 0BB= H IU   EJ  CA. ACB=   EU\n",
	" 1A-H V A=  CEU 0AEH 1AEN    E  CLU  ABG H IH A A= J B. A  9                    \n",
}

var egHelloWorld = "ABCDE" +
	"6" +
	"3000" + //                ORIG 3000
	"0000001187" + // START    IOC  0(18)
	"0787219621" + //          OUT  *+2(18)
	"0000000133" + //          HLT
	"0135582544" + //          ALF  HELLO
	"0006882509" + //          ALF   WORL
	"0067108864" + //          ALF  D
	"          \n" //          END  START

var egTransfer = "TRANS03000                                                                      \n"

var okLoader = "\014HELLO WORLD\n"
