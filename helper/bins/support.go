package bins

import (
	"fmt"

	"github.com/dosquad/mage/helper/paths"
	"github.com/dosquad/mage/loga"
	"github.com/princjef/mageutil/bintool"
)

//nolint:gochecknoglobals // ignore globals
var supportCmds map[string]*bintool.BinTool

func Support(name string) (*bintool.BinTool, error) {
	if supportCmds == nil {
		supportCmds = make(map[string]*bintool.BinTool)
	}

	switch name { //nolint:gocritic // ignore single case switch, for future support tools
	case "stringer":
		if supportCmds["stringer"] == nil {
			_ = Verdump().Ensure()
			ver := MustVersionLoadCache().GetVersion("golang.org/x/tools/cmd/stringer")
			loga.PrintInfof("stringer Version: %s", ver)
			supportCmds["stringer"] = bintool.Must(bintool.NewGo(
				"golang.org/x/tools/cmd/stringer",
				ver,
				bintool.WithFolder(paths.MustGetGoBin()),
				bintool.WithVersionCmd(paths.MustGetGoBin("verdump")+" mod {{.FullCmd}}"),
			))
		}

		return supportCmds["stringer"], nil
	}

	return nil, fmt.Errorf("unsupported tool: %s", name)
}
