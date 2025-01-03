package mix

import (
	"math"
	"testing"
)

const epsilon = 1e-5

func TestFixedToFloat(t *testing.T) {
	for i, ff := range egFixedFloat {
		res, ov := FixedToFloat(ff.fixed)
		if ov || int32(res) != int32(ff.float) {
			t.Errorf("#%d: got: %#v, %v (%v), want: %#v, %v",
				i+1, res, ov, res, ff.float, false)
		}
	}
}

func TestFloatToFixed(t *testing.T) {
	for i, ff := range egFixedFloat {
		res, ov := FloatToFixed(ff.float)
		if ov || int32(res) != int32(ff.fixed) {
			t.Errorf("#%d: got: %#v, %v (%v), want: %#v, %v",
				i+1, res, ov, res, ff.fixed, false)
		}
	}
}

func TestNativeToFloat(t *testing.T) {
	for i, nf := range egNativeFloat {
		res, st := NativeToFloat(nf.native)
		if st != nf.status ||
			(nf.status == 0 && int32(res) != int32(nf.float)) {
			t.Errorf("#%d: got: %#v, %v (%v), want: %#v, %v",
				i+1, res, st, res, nf.float, nf.status)
		}
	}
}

func TestFloatToNative(t *testing.T) {
	for i, nf := range egNativeFloat {
		if nf.status != 0 {
			continue
		}
		nat, _ := FloatToNative(nf.float)
		if math.Abs(nat-nf.native) > epsilon {
			t.Errorf("#%d: got: %#v, want: %#v",
				i+1, nat, nf.native)
		}
	}
}

var (
	egFixedFloat = []struct{ fixed, float Word }{
		{NewWord(1), NewWord(04101000000)},            // #1
		{NewWord(-511), NewWord(-04207770000)},        // #2
		{NewWord(-2097153), NewWord(-04410000001)},    // #3
		{NewWord(MaxWord - 63), NewWord(04577777777)}, // #4
	}

	egNativeFloat = []struct {
		native float64
		float  Word
		status int
	}{
		{1, NewWord(04101000000), 0},                         // #1
		{64, NewWord(04201000000), 0},                        // #2
		{1.0 / 64, NewWord(04001000000), 0},                  // #3
		{-4096, NewWord(-04301000000), 0},                    // #4
		{math.Pi, NewWord(04103110376), 0},                   // #5
		{9.807970876941034e+55, NewWord(MaxWord), 0},         // #6
		{math.MaxFloat64, NewWord(01320000000), expOverflow}, // #7
	}
)
