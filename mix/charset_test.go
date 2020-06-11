package mix

import (
	"reflect"
	"testing"
)

func TestConvertToUTF8(t *testing.T) {
	got := ConvertToUTF8(egMIX)
	if got != egUTF8 {
		t.Errorf(`got: %q, want: %q`, got, egUTF8)
	}
}

func TestConvertToMIX(t *testing.T) {
	got, err := ConvertToMIX(egUTF8)
	if err != nil {
		t.Errorf("error: %v", err)
	} else if !reflect.DeepEqual(got, egMIX) {
		t.Errorf("got: %#v, want: %#v", got, egMIX)
	}
}

var (
	egUTF8 = "HELLOWΦRLD"
	egMIX  = []Word{
		NewWord(01005151520),
		NewWord(03224231504),
	}
)
