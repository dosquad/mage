package helper

import (
	"runtime"

	"github.com/princjef/mageutil/bintool"
)

const (
	golangciLintVersion     = "1.59.0"
	govulncheckVersion      = "latest"
	protocGenGoGRPCVersion  = "latest"
	protocGenGoTwirpVersion = "v8.1.3"
)

//nolint:gochecknoglobals // ignore globals
var golangciLint *bintool.BinTool

func BinGolangCILint() *bintool.BinTool {
	if golangciLint == nil {
		ver := GetEnv("GOLANGCILINT_VERSION", golangciLintVersion)
		PrintInfo("Golang CI Lint Version: %s", ver)
		golangciLint = bintool.Must(bintool.New(
			"golangci-lint{{.BinExt}}",
			golangciLintVersion,
			"https://github.com/golangci/golangci-lint/releases/download/"+
				"v{{.Version}}/golangci-lint-{{.Version}}-{{.GOOS}}-{{.GOARCH}}{{.ArchiveExt}}",
			bintool.WithFolder(MustGetGoBin()),
		))
	}

	return golangciLint
}

//nolint:gochecknoglobals // ignore globals
var govulncheck *bintool.BinTool

func BinGovulncheck() *bintool.BinTool {
	if govulncheck == nil {
		ver := GetEnv("GOVULNCHECK_VERSION", govulncheckVersion)
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
var protoc *bintool.BinTool

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

		protocVer := GetProtocVersion()

		PrintInfo("Protocol Buffer Version: %s", protocVer)
		PanicIfError(ExtractArchive(
			"https://github.com/protocolbuffers/protobuf/releases/download/v"+protocVer+"/"+
				"protoc-"+protocVer+"-"+goOperatingSystem+"-"+goArch+".zip",
			MustGetWD("artifacts", "protobuf"),
		), "Extract Archive")

		protoc = bintool.Must(bintool.New(
			"protoc{{.BinExt}}",
			protocVer,
			"https://github.com/protocolbuffers/protobuf/releases/download/v"+protocVer+"/"+
				"protoc-"+protocVer+"-"+goOperatingSystem+"-"+goArch+".zip",
			bintool.WithFolder(MustGetWD("artifacts", "protobuf", "bin")),
		))
	}

	return protoc
}

//nolint:gochecknoglobals // ignore globals
var protocGenGo *bintool.BinTool

func BinProtocGenGo() *bintool.BinTool {
	if protocGenGo == nil {
		ver := GetEnv("PROTOCGENGO_VERSION", GetProtobufVersion())
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

func BinProtocGenGoGRPC() *bintool.BinTool {
	if protocGenGoGRPC == nil {
		ver := GetEnv("PROTOCGENGOGRPC_VERSION", protocGenGoGRPCVersion)
		PrintInfo("Protocol Buffer Golang gRPC Version: %s", ver)
		protocGenGoGRPC = bintool.Must(bintool.NewGo(
			"google.golang.org/grpc/cmd/protoc-gen-go-grpc",
			ver,
			bintool.WithFolder(MustGetProtobufPath()),
		))
	}

	return protocGenGoGRPC
}

//nolint:gochecknoglobals // ignore globals
var protocGenGoTwirp *bintool.BinTool

func BinProtocGenGoTwirp() *bintool.BinTool {
	if protocGenGoTwirp == nil {
		ver := GetEnv("PROTOCGENGOTWIRP_VERSION", protocGenGoTwirpVersion)
		PrintInfo("Protocol Buffer Golang Twirp Version: %s", ver)
		protocGenGoTwirp = bintool.Must(bintool.NewGo(
			"github.com/twitchtv/twirp/protoc-gen-twirp",
			ver,
			bintool.WithFolder(MustGetProtobufPath()),
		))
	}

	return protocGenGoTwirp
}
