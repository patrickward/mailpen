package mailpen_test

import (
	"io/fs"
	"os"
	"testing"
)

// testFS creates a testing filesystem from the testdata directory
func testFS(t *testing.T, dir string) fs.FS {
	t.Helper()
	return os.DirFS("testdata/" + dir)
}
