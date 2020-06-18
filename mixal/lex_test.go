package mixal

import (
	"testing"

	"github.com/jflude/gnuth/mix"
)

func TestLexical(t *testing.T) {
	asmbLex.symbols = make(map[string]mix.Word)
	asmbLex.symbols["L"] = mix.NewWord(2000)
	asmbLex.fixups = make(map[string][]int)
	for i := range egMatch {
		asmbLex.self = 3000
		asmbLex.input = egMatch[i].input
		asmbLex.tokens = nil

		ok := egMatch[i].lexer()
		if !ok || asmbLex.input != okMatch[i].remain {
			t.Errorf(`%d: got: %q, %v, want: %q, true`,
				i+1, asmbLex.input, ok, okMatch[i].remain)
		}
		if len(asmbLex.tokens) == 0 {
			t.Errorf("%d: got: nil, want: %v",
				i+1, okMatch[i].tok)
			continue
		}
		for j, v := range okMatch[i].tok {
			if j >= len(asmbLex.tokens) {
				t.Errorf("%d: got: nil, want: %v", i+1, v)
				break
			}
			if v != asmbLex.tokens[j] {
				t.Errorf("%d: got: %v, want: %v",
					i+1, asmbLex.tokens[j], v)
				break
			}
		}
	}
}

var asmbLex asmb

var egMatch = []struct {
	input string
	lexer func() bool
}{
	{"BUF1 ", asmbLex.matchFutureRef},              // #1
	{"12345678900 ", asmbLex.matchNumber},          // #2
	{"JAZ", asmbLex.matchOperator},                 // #3
	{"* ", asmbLex.matchAsterisk},                  // #4
	{"*** X", asmbLex.parseExpr},                   // #5
	{"-1+5*20//6+* Y", asmbLex.parseExpr},          // #6
	{"(1:5) Z", asmbLex.parseFPart},                // #7
	{"1000(1:3),2000(4:5) X", asmbLex.parseWValue}, // #8
	{"=1-L= Y", asmbLex.parseLiteral},              // #9
}

var okMatch = []struct {
	remain string
	tok    []token
}{
	{" ", []token{{symbol, "BUF1"}}},      // #1
	{"0 ", []token{{number, 1234567890}}}, // #2
	{"", []token{{operator, "JAZ"}}},      // #3
	{" ", []token{{asterisk, 3000}}},      // #4
	{" X", []token{ // #5
		{asterisk, 3000},
		{'*', nil},
		{asterisk, 3000}}},
	{" Y", []token{ // #6
		{'-', nil},
		{number, 1},
		{'+', nil},
		{number, 5},
		{'*', nil},
		{number, 20},
		{'\\', nil},
		{number, 6},
		{'+', nil},
		{asterisk, 3000}}},
	{" Z", []token{ // #7
		{'(', nil},
		{number, 1},
		{':', nil},
		{number, 5},
		{')', nil}}},
	{" X", []token{ // #8
		{number, 1000},
		{'(', nil},
		{number, 1},
		{':', nil},
		{number, 3},
		{')', nil},
		{',', nil},
		{number, 2000},
		{'(', nil},
		{number, 4},
		{':', nil},
		{number, 5},
		{')', nil}}},
	{" Y", []token{ // #9
		{'=', nil},
		{number, 1},
		{'-', nil},
		{symbol, "L"},
		{'=', nil}}},
}
