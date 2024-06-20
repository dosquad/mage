package helper

import (
	"runtime"

	"github.com/princjef/mageutil/bintool"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//nolint:gochecknoglobals // ignore globals
var bufTool *bintool.BinTool

func Buf() *bintool.BinTool {
	if bufTool == nil {
		ver := MustVersionLoadCache().GetVersion(BufVersion)
		PrintInfo("Buf Version: %s", ver)
		bufTool = bintool.Must(bintool.NewGo(
			"github.com/bufbuild/buf/cmd/buf",
			ver,
			bintool.WithFolder(MustGetGoBin()),
		))
	}

	return bufTool
}

//nolint:gochecknoglobals // ignore globals
var golangciLint *bintool.BinTool

// BinGolangCILint returns a singleton for golangci-lint.
func BinGolangCILint() *bintool.BinTool {
	if golangciLint == nil {
		// ver := GetEnv("GOLANGCILINT_VERSION", golangciLintVersion)
		ver := MustVersionLoadCache().GetVersion(GolangciLintVersion)
		PrintInfo("Golang CI Lint Version: %s", ver)
		golangciLint = bintool.Must(bintool.New(
			"golangci-lint{{.BinExt}}",
			ver,
			"https://github.com/golangci/golangci-lint/releases/download/"+
				"v{{.Version}}/golangci-lint-{{.Version}}-{{.GOOS}}-{{.GOARCH}}{{.ArchiveExt}}",
			bintool.WithFolder(MustGetGoBin()),
		))
	}

	return golangciLint
}

//nolint:gochecknoglobals // ignore globals
var govulncheck *bintool.BinTool

// BinGovulncheck returns a singleton for govulncheck.
func BinGovulncheck() *bintool.BinTool {
	if govulncheck == nil {
		// ver := GetEnv("GOVULNCHECK_VERSION", govulncheckVersion)
		ver := MustVersionLoadCache().GetVersion(GovulncheckVersion)
		PrintInfo("Golang Vulnerability Check Version: %s", ver)
		govulncheck = bintool.Must(bintool.NewGo(
			"golang.org/x/vuln/cmd/govulncheck",
			ver,
			bintool.WithFolder(MustGetGoBin()),
			bintool.WithVersionCmd("{{.FullCmd}} -version"),
		))
	}

	return govulncheck
}

//nolint:gochecknoglobals // ignore globals
var goreleaser *bintool.BinTool

// BinGoreleaser returns a singleton for goreleaser.
func BinGoreleaser() *bintool.BinTool {
	if goreleaser == nil {
		ver := MustVersionLoadCache().GetVersion(GoreleaserVersion)
		PrintInfo("Goreleaser Version: %s", ver)

		goOperatingSystem, goArch := runtime.GOOS, runtime.GOARCH
		goOperatingSystem = cases.Title(language.English).String(goOperatingSystem)

		if runtime.GOARCH == "amd64" {
			goArch = "x86_64"
		}

		url := "https://github.com/goreleaser/goreleaser/releases/download/v" + ver + "/" +
			"goreleaser_" + goOperatingSystem + "_" + goArch + "{{.ArchiveExt}}"
		PrintDebug("Goreleaser URL: %s", url)

		goreleaser = bintool.Must(bintool.New(
			"goreleaser{{.BinExt}}",
			ver,
			url,
			bintool.WithFolder(MustGetGoBin()),
		))
	}

	return goreleaser
}

//nolint:gochecknoglobals // ignore globals
var protoc *bintool.BinTool

// BinProtoc returns a singleton for protoc, also downloads the includes.
func BinProtoc() *bintool.BinTool {
	if protoc == nil {
		goOperatingSystem, goArch := runtime.GOOS, runtime.GOARCH
		if runtime.GOOS == "darwin" {
			goOperatingSystem = "osx"
		}

		switch runtime.GOARCH {
		case "amd64":
			goArch = "x86_64"
		case "arm64":
			goArch = "aarch_64"
		}

		protocVer := MustVersionLoadCache().GetVersion(ProtocVersion)

		PrintInfo("Protocol Buffer Version: %s", protocVer)
		PanicIfError(ExtractArchive(
			"https://github.com/protocolbuffers/protobuf/releases/download/v"+protocVer+"/"+
				"protoc-"+protocVer+"-"+goOperatingSystem+"-"+goArch+".zip",
			MustGetArtifactPath("protobuf"),
		), "Extract Archive")

		protoc = bintool.Must(bintool.New(
			"protoc{{.BinExt}}",
			protocVer,
			"https://github.com/protocolbuffers/protobuf/releases/download/v"+protocVer+"/"+
				"protoc-"+protocVer+"-"+goOperatingSystem+"-"+goArch+".zip",
			bintool.WithFolder(MustGetArtifactPath("protobuf", "bin")),
		))
	}

	return protoc
}

//nolint:gochecknoglobals // ignore globals
var protocGenGo *bintool.BinTool

// BinProtocGenGo returns a singleton for protoc-gen-go.
func BinProtocGenGo() *bintool.BinTool {
	if protocGenGo == nil {
		// ver := GetEnv("PROTOCGENGO_VERSION", GetProtobufVersion())
		ver := MustVersionLoadCache().GetVersion(ProtocGenGoVersion)
		PrintInfo("Protocol Buffer Golang Version: %s", ver)
		protocGenGo = bintool.Must(bintool.NewGo(
			"google.golang.org/protobuf/cmd/protoc-gen-go",
			ver,
			bintool.WithFolder(MustGetProtobufPath()),
		))
	}

	return protocGenGo
}

//nolint:gochecknoglobals // ignore globals
var protocGenGoGRPC *bintool.BinTool

// BinProtocGenGoGRPC returns a singleton for protoc-gen-go-grpc.
func BinProtocGenGoGRPC() *bintool.BinTool {
	if protocGenGoGRPC == nil {
		_ = BinVerdump().Ensure()
		ver := MustVersionLoadCache().GetVersion(ProtocGenGoGRPCVersion)
		PrintInfo("Protocol Buffer Golang gRPC Version: %s", ver)
		protocGenGoGRPC = bintool.Must(bintool.NewGo(
			"google.golang.org/grpc/cmd/protoc-gen-go-grpc",
			ver,
			bintool.WithFolder(MustGetProtobufPath()),
			bintool.WithVersionCmd(MustGetGoBin("verdump")+" mod {{.FullCmd}}"),
		))
	}

	return protocGenGoGRPC
}

//nolint:gochecknoglobals // ignore globals
var protocGenGoTwirp *bintool.BinTool

// BinProtocGenGoTwirp returns a singleton for protoc-gen-twirp.
func BinProtocGenGoTwirp() *bintool.BinTool {
	if protocGenGoTwirp == nil {
		// ver := GetEnv("PROTOCGENGOTWIRP_VERSION", protocGenGoTwirpVersion)
		ver := MustVersionLoadCache().GetVersion(ProtocGenGoTwirpVersion)
		PrintInfo("Protocol Buffer Golang Twirp Version: %s", ver)
		protocGenGoTwirp = bintool.Must(bintool.NewGo(
			"github.com/twitchtv/twirp/protoc-gen-twirp",
			ver,
			bintool.WithFolder(MustGetProtobufPath()),
		))
	}

	return protocGenGoTwirp
}

//nolint:gochecknoglobals // ignore globals
var yq *bintool.BinTool

// BinYQ returns a singleton for yq.
func BinYQ() *bintool.BinTool {
	if yq == nil {
		ver := MustVersionLoadCache().GetVersion(YQVersion)
		PrintInfo("YQ Version: %s", ver)
		yq = bintool.Must(bintool.NewGo(
			"github.com/mikefarah/yq/v4",
			ver,
			bintool.WithFolder(MustGetGoBin()),
		))
	}

	return yq
}

//nolint:gochecknoglobals // ignore globals
var wirebin *bintool.BinTool

// BinWire returns a singleton for wirebin.
func BinWire() *bintool.BinTool {
	if wirebin == nil {
		_ = BinVerdump().Ensure()
		ver := MustVersionLoadCache().GetVersion(WireVersion)
		PrintInfo("Wire Version: %s", ver)
		wirebin = bintool.Must(bintool.NewGo(
			"github.com/google/wire/cmd/wire",
			ver,
			bintool.WithFolder(MustGetGoBin()),
			bintool.WithVersionCmd(MustGetGoBin("verdump")+" mod {{.FullCmd}}"),
		))
	}

	return wirebin
}

//nolint:gochecknoglobals // ignore globals
var verdump *bintool.BinTool

// BinVerdump returns a singleton for verdump.
func BinVerdump() *bintool.BinTool {
	if verdump == nil {
		ver := MustVersionLoadCache().GetVersion(VerdumpVersion)
		PrintInfo("Verdump Version: %s", ver)
		verdump = bintool.Must(bintool.NewGo(
			"github.com/dosquad/mage/cmd/verdump",
			ver,
			bintool.WithFolder(MustGetGoBin()),
		))
	}

	return verdump
}
