package clipboard

import (
	"os/exec"
	"runtime"
	"testing"
)

func TestWrite(t *testing.T) {
	if runtime.GOOS == "darwin" {
		if _, err := exec.LookPath("pbcopy"); err != nil {
			t.Skip("pbcopy not available")
		}
	} else if runtime.GOOS == "linux" {
		if _, err := exec.LookPath("xclip"); err != nil {
			if _, err := exec.LookPath("xsel"); err != nil {
				t.Skip("no clipboard tool available")
			}
		}
	} else if runtime.GOOS == "windows" {
		if _, err := exec.LookPath("clip"); err != nil {
			t.Skip("clip not available")
		}
	} else {
		t.Skipf("unsupported platform: %s", runtime.GOOS)
	}

	err := Write("gh-compare test")
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
}
