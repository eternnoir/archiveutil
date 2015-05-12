package archiveutil

import (
	"bytes"
	"testing"
)

func TestZipAddFile(t *testing.T) {
	buf := new(bytes.Buffer)
	z := CreateArchive("zip", buf)
	err := z.AddFile("./README.md")
	if err != nil {
		t.Fail()
	}
	err = z.Close()
	if err != nil {
		t.Fail()
	}
}
