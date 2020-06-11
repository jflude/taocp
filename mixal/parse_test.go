package mixal

import (
	"strings"
	"testing"
)

func TestParseLine(t *testing.T) {
	r := strings.NewReader(egReadProgram)
	var a asmb
	if err := a.readProgram(r, func(a *asmb, loc, op, address string) {
		t.Helper()
		t.Logf("parse: %s %s %s", loc, op, address)
		parseLine(a, loc, op, address)
		t.Logf("%d: %v", a.count, a.tokens)
	}); err != nil {
		t.Error("error:", err)
	}
}
