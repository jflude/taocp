package mix

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestGoButton(t *testing.T) {
	c, tmpDir := newSandbox(t, "")
	defer closeSandbox(t, c, tmpDir)
	if testing.Verbose() {
		c.Tracer = os.Stdout
		c.Trigger = 0
	}

	deck := append(append([]byte(nil), egLoader[0]...), egLoader[1]...)
	deck = append(append(deck, egHelloWorld...), egTransfer...)
	err := os.WriteFile("reader.mix", deck, 0644)
	if err != nil {
		t.Fatal("error:", err)
	}
	if err = c.GoButton(); !errors.Is(err, ErrHalted) {
		t.Fatal("error:", err)
	}
	b, err := os.ReadFile("printer.mix")
	if err != nil {
		t.Fatal("error:", err)
	}
	if strings.Compare(string(b), okLoader) != 0 {
		t.Error("got: incorrect printer output")
	}
}

var egLoader = []string{
	//   5   10   15   20   25   30   35   40   45   50   55   60   65   70   75   80
	" O O6 Z O6    I C O4 0 EH A  F F CF 0  E   EU 0 IH G BB   EJ  CA. Z EU   EH E BA\n",
	"   EU 2A-H S BB  C U 1AEH 2AEN V  E  CLU  ABG Z EH E BB J B. A  9               \n",
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
