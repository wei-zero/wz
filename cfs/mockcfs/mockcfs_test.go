package mockcfs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMockcfs(t *testing.T) {
	cfs, err := New()
	if err != nil {
		t.Error(err)
		return
	}

	err = cfs.PutObject("bucket", "key", "mockcfs.go", "content-type")
	if err != nil {
		t.Error(err)
		return
	}

	dstFile := filepath.Join(os.TempDir(), "mockcfs.go")
	err = cfs.GetObject("bucket", "key", dstFile)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("mockcfs.go is downloaded to", dstFile)
}
