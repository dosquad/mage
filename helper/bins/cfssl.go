package bins

import (
	"github.com/dosquad/mage/helper/paths"
	"github.com/dosquad/mage/loga"
	"github.com/princjef/mageutil/bintool"
)

//nolint:gochecknoglobals // ignore globals
var cfsslCmd *bintool.BinTool

// Cfssl returns a singleton for cfssl.
func Cfssl() *bintool.BinTool {
	if cfsslCmd == nil {
		_ = Verdump().Ensure()
		ver := MustVersionLoadCache().GetVersion(CFSSLVersion)
		loga.PrintInfo("cfssl Version: %s", ver)
		cfsslCmd = bintool.Must(bintool.NewGo(
			"github.com/cloudflare/cfssl/cmd/cfssl",
			ver,
			bintool.WithFolder(paths.MustGetGoBin()),
			bintool.WithVersionCmd(paths.MustGetGoBin("verdump")+" mod {{.FullCmd}}"),
		))
	}

	return cfsslCmd
}

//nolint:gochecknoglobals // ignore globals
var cfsslJSONCmd *bintool.BinTool

// CfsslJSON returns a singleton for cfssl.
func CfsslJSON() *bintool.BinTool {
	if cfsslJSONCmd == nil {
		_ = Verdump().Ensure()
		ver := MustVersionLoadCache().GetVersion(CFSSLVersion)
		loga.PrintInfo("cfssl Version: %s", ver)
		cfsslJSONCmd = bintool.Must(bintool.NewGo(
			"github.com/cloudflare/cfssl/cmd/cfssljson",
			ver,
			bintool.WithFolder(paths.MustGetGoBin()),
			bintool.WithVersionCmd(paths.MustGetGoBin("verdump")+" mod {{.FullCmd}}"),
		))
	}

	return cfsslJSONCmd
}
