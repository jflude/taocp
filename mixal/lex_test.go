package mixal

import "testing"

func TestMatching(t *testing.T) {
	asmbMatch.symbols = make(map[string]int)
	asmbMatch.symbols["L"] = 2000
	for i := 0; i < len(egMatch); i++ {
		asmbMatch.input = egMatch[i].input
		asmbMatch.tokens = nil
		asmbMatch.self = 3000
		ok := egMatch[i].lexer()
		if !ok || asmbMatch.input != okMatch[i].remain {
			t.Errorf(`%d: got: %q, %v, want: %q, true`,
				i+1, asmbMatch.input, ok, okMatch[i].remain)
		}
		if len(asmbMatch.tokens) == 0 {
			t.Errorf("%d: got: nil, want: %v",
				i+1, okMatch[i].tok)
			continue
		}
		for j, v := range okMatch[i].tok {
			if j >= len(asmbMatch.tokens) {
				t.Errorf("%d: got: nil, want: %v", i+1, v)
				break
			}
			if v != asmbMatch.tokens[j] {
				t.Errorf("%d: got: %v, want: %v",
					i+1, asmbMatch.tokens[j], v)
				break
			}
		}
	}
}

var asmbMatch asmb

var egMatch = []struct {
	input string
	lexer func() bool
}{
	{"BUF1 ", asmbMatch.matchSymbol},
	{"12345678900 ", asmbMatch.matchNumber},
	{"JAZ", asmbMatch.matchOperator},
	{"* ", asmbMatch.matchAsterisk},
	{"*** X", asmbMatch.matchExpr},
	{"-1+5*20//6+* Y", asmbMatch.matchExpr},
	{"(1:5) Z", asmbMatch.matchFPart},
	{"1000(1:3),2000(4:5) X", asmbMatch.matchWValue},
	{"=1-L= Y", asmbMatch.matchLiteral},
}

var okMatch = []struct {
	remain string
	tok    []token
}{
	{" ", []token{{symbol, "BUF1"}}},
	{"0 ", []token{{number, 1234567890}}},
	{"", []token{{operator, "JAZ"}}},
	{" ", []token{{asterisk, 3000}}},
	{" X", []token{
		{asterisk, 3000},
		{'*', nil},
		{asterisk, 3000}}},
	{" Y", []token{
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
	{" Z", []token{
		{'(', nil},
		{number, 1},
		{':', nil},
		{number, 5},
		{')', nil}}},
	{" X", []token{
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
	{" Y", []token{
		{'=', nil},
		{number, 1},
		{'-', nil},
		{symbol, "L"},
		{'=', nil}}},
}
