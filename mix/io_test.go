package mix

import "testing"

func TestIO(t *testing.T) {
	c, tmpDir := newSandbox(t, "")
	defer closeSandbox(t, c, tmpDir)

	// DRUM10
	if err := c.bindDevice(10); err != nil {
		t.Fatal("error:", err)
	}
	drum := c.Devices[10]
	buf := make([]Word, drum.BlockSize())
	for i := 0; i < 10; i++ {
		for j := range buf {
			buf[j] = NewWord(100*i + j)
		}
		c.Reg[X] = NewWord(10 * i)
		if _, err := drum.Write(buf); err != nil {
			t.Error("error:", err)
		}
	}
	for i := 0; i < 10; i++ {
		c.Reg[X] = NewWord(10 * i)
		if _, err := drum.Read(buf); err != nil {
			t.Error("error:", err)
		}
		for j, v := range buf {
			want := 100*i + j
			if v.Int() != want {
				t.Errorf("#%d.%d: got: %d, want: %d",
					i, j, v.Int(), want)
				break
			}
		}
	}

	// TAPE02
	if err := c.bindDevice(2); err != nil {
		t.Fatal("error:", err)
	}
	tape := c.Devices[2]
	buf = make([]Word, tape.BlockSize())
	for i := 0; i < 10; i++ {
		for j := range buf {
			buf[j] = NewWord(100*i + j)
		}
		if _, err := tape.Write(buf); err != nil {
			t.Error("error:", err)
		}
	}
	if _, err := tape.Control(0); err != nil {
		t.Error("error:", err)
	}
	if _, err := tape.Control(1); err != nil {
		t.Error("error:", err)
	}
	if _, err := tape.Control(1); err != nil {
		t.Error("error:", err)
	}
	for i := 2; i < 10; i++ {
		if _, err := tape.Read(buf); err != nil {
			t.Error("error:", err)
		}
		for j, v := range buf {
			want := 100*i + j
			if v.Int() != want {
				t.Errorf("got: %d, want: %d",
					v.Int(), want)
				break
			}
		}
	}
}
