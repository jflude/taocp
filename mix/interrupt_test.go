package mix

import "testing"

func TestInterrupt(t *testing.T) {
	c, tmpDir := newSandbox(t, "")
	defer closeSandbox(t, c, tmpDir)
	for i := 0; i < len(c.Reg); i++ {
		c.Reg[i] = NewWord(10 * (i + 1))
	}
	c.Contents[mBase-10] = NewWord(10)
	printIntr(t, c, 500)
	printIntr(t, c, 1000)
	printIntr(t, c, 5001)
	printIntr(t, c, 20002)
	if _, err := c.ioc(0, 0, 18, IOC, 0); err != nil {
		t.Error("error:", err)
	}
	c.ctrl = false
	printIntr(t, c, 35000)
	printIntr(t, c, 100000)
	printIntr(t, c, 1020002)
	printIntr(t, c, 1020004)
}

func printIntr(t *testing.T, c *Computer, elapsed int64) {
	c.Elapsed = elapsed
	c.interrupt(c.ctrl)
	t.Log("next:", c.next, "ctrl:", c.ctrl, "Elapsed:", c.Elapsed,
		"lastTick:", c.lastTick, "pending:", c.pending,
		"Contents:", c.Contents[mBase-10:mBase])
}
