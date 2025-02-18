package paths_test

import (
	"os"
	"testing"

	"github.com/dosquad/mage/helper/paths"
	"github.com/na4ma4/go-permbits"
)

func TestFileExistsInPath(t *testing.T) {
	if v := paths.FileExistsInPath("*.randomfile", "artifacts"); v {
		t.Errorf("FileExistsInPath: [artifacts](*.randomfile) got '%t', want '%t'", v, false)
	}

	localPath := paths.MustGetArtifactPath("testdata")

	_ = os.MkdirAll(localPath, permbits.MustString("ug=rwx,o=rx"))
	_ = os.WriteFile(localPath+"/test.proto", []byte("testfile"), permbits.MustString("ug=rw,o=r"))

	if v := paths.FileExistsInPath("*.proto", paths.MustGetArtifactPath()); !v {
		t.Errorf("FileExistsInPath: [artifacts/testdata](*.proto) got '%t', want '%t'", v, true)
	}
}
