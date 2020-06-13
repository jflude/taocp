package mix

import (
	"io/ioutil"
	"os"
	"testing"
)

func newSandbox(t *testing.T, tmpDir string) (*Computer, string) {
	var ret string
	var err error
	if tmpDir == "" {
		tmpDir, err = ioutil.TempDir("", "gnuth-mix-test")
		if err != nil {
			t.Fatal("error:", err)
		}
		ret = tmpDir
	}
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatal("error:", err)
	}
	return NewComputer(nil), ret
}

func closeSandbox(t *testing.T, c *Computer, tmpDir string) {
	if err := c.Shutdown(); err != nil {
		t.Error("error:", err)
	}
	if tmpDir != "" {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Error("error:", err)
		}
	}
}
