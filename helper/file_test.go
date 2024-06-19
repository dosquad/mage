package helper_test

import (
	"os"
	"testing"

	"github.com/dosquad/mage/helper"
	"github.com/na4ma4/go-permbits"
)

func TestFileExistsInPath(t *testing.T) {
	if v := helper.FileExistsInPath("*.randomfile", "artifacts"); v {
		t.Errorf("FileExistsInPath: [artifacts](*.randomfile) got '%t', want '%t'", v, false)
	}

	_ = os.MkdirAll("artifacts/testdata", permbits.MustString("ug=rwx,o=rx"))
	_ = os.WriteFile("artifacts/testdata/test.proto", []byte("testfile"), permbits.MustString("ug=rw,o=r"))

	if v := helper.FileExistsInPath("*.proto", "artifacts"); !v {
		t.Errorf("FileExistsInPath: [artifacts/testdata](*.proto) got '%t', want '%t'", v, true)
	}
}
