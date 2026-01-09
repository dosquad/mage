package bins

import (
	"fmt"
	"runtime"

	"github.com/dosquad/mage/helper/must"
	"github.com/dosquad/mage/helper/paths"
	"github.com/dosquad/mage/helper/web"
	"github.com/dosquad/mage/loga"
	"github.com/princjef/mageutil/bintool"
	"github.com/princjef/mageutil/shellcmd"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//nolint:gochecknoglobals // ignore globals
var bufTool *bintool.BinTool

func Buf() *bintool.BinTool {
	if bufTool == nil {
		ver := MustVersionLoadCache().GetVersion(BufVersion)
		loga.PrintInfof("Buf Version: %s", ver)
		bufTool = bintool.Must(bintool.NewGo(
			"github.com/bufbuild/buf/cmd/buf",
			ver,
			bintool.WithFolder(paths.MustGetGoBin()),
		))
	}

	return bufTool
}

//nolint:gochecknoglobals // ignore globals
var golangciLint *bintool.BinTool

// GolangCILint returns a singleton for golangci-lint.
func GolangCILint() *bintool.BinTool {
	if golangciLint == nil {
		// ver := GetEnv("GOLANGCILINT_VERSION", golangciLintVersion)
		ver := MustVersionLoadCache().GetVersion(GolangciLintVersion)
		loga.PrintInfof("Golang CI Lint Version: %s", ver)
		golangciLint = bintool.Must(bintool.New(
			"golangci-lint{{.BinExt}}",
			ver,
			"https://github.com/golangci/golangci-lint/releases/download/"+
				"v{{.Version}}/golangci-lint-{{.Version}}-{{.GOOS}}-{{.GOARCH}}{{.ArchiveExt}}",
			bintool.WithFolder(paths.MustGetGoBin()),
		))
	}

	return golangciLint
}

//nolint:gochecknoglobals // ignore globals
var govulncheck *bintool.BinTool

// Govulncheck returns a singleton for govulncheck.
func Govulncheck() *bintool.BinTool {
	if govulncheck == nil {
		// ver := GetEnv("GOVULNCHECK_VERSION", govulncheckVersion)
		ver := MustVersionLoadCache().GetVersion(GovulncheckVersion)
		loga.PrintInfof("Golang Vulnerability Check Version: %s", ver)
		govulncheck = bintool.Must(bintool.NewGo(
			"golang.org/x/vuln/cmd/govulncheck",
			ver,
			bintool.WithFolder(paths.MustGetGoBin()),
			bintool.WithVersionCmd("{{.FullCmd}} -version"),
		))
	}

	return govulncheck
}

//nolint:gochecknoglobals // ignore globals
var goreleaser *bintool.BinTool

// Goreleaser returns a singleton for goreleaser.
func Goreleaser() *bintool.BinTool {
	if goreleaser == nil {
		ver := MustVersionLoadCache().GetVersion(GoreleaserVersion)
		loga.PrintInfof("Goreleaser Version: %s", ver)

		goOperatingSystem, goArch := runtime.GOOS, runtime.GOARCH
		goOperatingSystem = cases.Title(language.English).String(goOperatingSystem)

		if runtime.GOARCH == "amd64" {
			goArch = "x86_64"
		}

		url := "https://github.com/goreleaser/goreleaser/releases/download/v" + ver + "/" +
			"goreleaser_" + goOperatingSystem + "_" + goArch + "{{.ArchiveExt}}"
		loga.PrintDebugf("Goreleaser URL: %s", url)

		goreleaser = bintool.Must(bintool.New(
			"goreleaser{{.BinExt}}",
			ver,
			url,
			bintool.WithFolder(paths.MustGetGoBin()),
		))
	}

	return goreleaser
}

//nolint:gochecknoglobals // ignore globals
var yq *bintool.BinTool

// YQ returns a singleton for yq.
func YQ() *bintool.BinTool {
	if yq == nil {
		ver := MustVersionLoadCache().GetVersion(YQVersion)
		loga.PrintInfof("YQ Version: %s", ver)
		yq = bintool.Must(bintool.NewGo(
			"github.com/mikefarah/yq/v4",
			ver,
			bintool.WithFolder(paths.MustGetGoBin()),
		))
	}

	return yq
}

//nolint:gochecknoglobals // ignore globals
var wirebin *bintool.BinTool

// Wire returns a singleton for wirebin.
func Wire() *bintool.BinTool {
	if wirebin == nil {
		_ = Verdump().Ensure()
		ver := MustVersionLoadCache().GetVersion(WireVersion)
		loga.PrintInfof("Wire Version: %s", ver)
		wirebin = bintool.Must(bintool.NewGo(
			"github.com/google/wire/cmd/wire",
			ver,
			bintool.WithFolder(paths.MustGetGoBin()),
			bintool.WithVersionCmd(paths.MustGetGoBin("verdump")+" mod {{.FullCmd}}"),
		))
	}

	return wirebin
}

//nolint:gochecknoglobals // ignore globals
var verdump *bintool.BinTool

// Verdump returns a singleton for verdump.
func Verdump() *bintool.BinTool {
	if verdump == nil {
		ver := MustVersionLoadCache().GetVersion(VerdumpVersion)
		loga.PrintInfof("Verdump Version: %s", ver)
		verdump = bintool.Must(bintool.NewGo(
			"github.com/dosquad/mage/cmd/verdump",
			ver,
			bintool.WithFolder(paths.MustGetGoBin()),
		))
	}

	return verdump
}

//nolint:gochecknoglobals // ignore globals
var kustomize *bintool.BinTool

func installKustomize() error {
	var instScript string
	{
		var err error
		instScript, err = web.DownloadToCache(
			"https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh",
		)
		if err != nil {
			return err
		}
	}

	ver := MustVersionLoadCache().GetVersion(KustomizeVersion)
	loga.PrintInfof("Kustomize Version: %s", ver)
	instCmd := shellcmd.Command(
		fmt.Sprintf(
			`bash "%s" "%s" "%s"`,
			instScript,
			ver,
			paths.MustGetArtifactPath("bin"),
		),
	)
	if err := instCmd.Run(); err != nil {
		return err
	}

	loga.PrintDebugf("Install script completed")

	return nil
}

// Kustomize returns a singleton for kustomize.
func Kustomize() *bintool.BinTool {
	if kustomize == nil {
		if !paths.FileExists(paths.MustGetArtifactPath("bin", "kustomize")) {
			must.PanicIfError(installKustomize(), "unable to install kustomize")
		}

		kustomize = bintool.Must(bintool.New("kustomize", "", "",
			bintool.WithFolder(paths.MustGetArtifactPath("bin")),
			bintool.WithVersionCmd(`{{.FullCmd}} version`),
		))
	}

	return kustomize
}

//nolint:gochecknoglobals // ignore globals
var kubeControllerGen *bintool.BinTool

// KubeControllerGen returns a singleton for kubeControllerGen.
func KubeControllerGen() *bintool.BinTool {
	if kubeControllerGen == nil {
		ver := MustVersionLoadCache().GetVersion(KubeControllerGenVersion)
		loga.PrintInfof("sigs.k8s.io Controller Gen Version: %s", ver)
		kubeControllerGen = bintool.Must(bintool.NewGo(
			"sigs.k8s.io/controller-tools/cmd/controller-gen",
			ver,
			bintool.WithFolder(paths.MustGetGoBin()),
		))
	}

	return kubeControllerGen
}

//nolint:gochecknoglobals // ignore globals
var kubeControllerEnvTest *bintool.BinTool

// KubeControllerEnvTest returns a singleton for kubeControllerEnvTest.
func KubeControllerEnvTest() *bintool.BinTool {
	if kubeControllerEnvTest == nil {
		ver := MustVersionLoadCache().GetVersion(KubeControllerEnvTestVersion)
		loga.PrintInfof("sigs.k8s.io Controller Runtime Version: %s", ver)
		kubeControllerEnvTest = bintool.Must(bintool.NewGo(
			"sigs.k8s.io/controller-runtime/tools/setup-envtest",
			ver,
			bintool.WithFolder(paths.MustGetGoBin()),
		))
	}

	return kubeControllerEnvTest
}

//nolint:gochecknoglobals // ignore globals
var vgtCmd *bintool.BinTool

// VGT returns a singleton for vgt.
func VGT() *bintool.BinTool {
	if vgtCmd == nil {
		_ = Verdump().Ensure()
		ver := MustVersionLoadCache().GetVersion(VGTVersion)
		loga.PrintInfof("vgt Version: %s", ver)
		vgtCmd = bintool.Must(bintool.NewGo(
			"github.com/roblaszczak/vgt",
			ver,
			bintool.WithFolder(paths.MustGetGoBin()),
			bintool.WithVersionCmd(paths.MustGetGoBin("verdump")+" mod {{.FullCmd}}"),
		))
	}

	return vgtCmd
}
