package lvmreport

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFromFileNotFound(t *testing.T) {
	path := filepath.Join(t.TempDir(), "does", "not", "exist")

	if _, err := FromFile(path); !os.IsNotExist(err) {
		t.Errorf("File should not be found: %v", err)
	}
}
