package bins

import (
	"runtime"

	"github.com/dosquad/mage/helper/must"
	"github.com/dosquad/mage/helper/paths"
	"github.com/dosquad/mage/helper/web"
	"github.com/dosquad/mage/loga"
	"github.com/princjef/mageutil/bintool"
)

//nolint:gochecknoglobals // ignore globals
var protoc *bintool.BinTool

// Protoc returns a singleton for protoc, also downloads the includes.
func Protoc() *bintool.BinTool {
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

		loga.PrintInfo("Protocol Buffer Version: %s", protocVer)
		must.PanicIfError(web.ExtractArchive(
			"https://github.com/protocolbuffers/protobuf/releases/download/v"+protocVer+"/"+
				"protoc-"+protocVer+"-"+goOperatingSystem+"-"+goArch+".zip",
			paths.MustGetArtifactPath("protobuf"),
		), "Extract Archive")

		protoc = bintool.Must(bintool.New(
			"protoc{{.BinExt}}",
			protocVer,
			"https://github.com/protocolbuffers/protobuf/releases/download/v"+protocVer+"/"+
				"protoc-"+protocVer+"-"+goOperatingSystem+"-"+goArch+".zip",
			bintool.WithFolder(paths.MustGetArtifactPath("protobuf", "bin")),
		))
	}

	return protoc
}

//nolint:gochecknoglobals // ignore globals
var protocGenGo *bintool.BinTool

// ProtocGenGo returns a singleton for protoc-gen-go.
func ProtocGenGo() *bintool.BinTool {
	if protocGenGo == nil {
		// ver := GetEnv("PROTOCGENGO_VERSION", GetProtobufVersion())
		ver := MustVersionLoadCache().GetVersion(ProtocGenGoVersion)
		loga.PrintInfo("Protocol Buffer Golang Version: %s", ver)
		protocGenGo = bintool.Must(bintool.NewGo(
			"google.golang.org/protobuf/cmd/protoc-gen-go",
			ver,
			bintool.WithFolder(paths.MustGetProtobufPath()),
		))
	}

	return protocGenGo
}

//nolint:gochecknoglobals // ignore globals
var protocGenGoGRPC *bintool.BinTool

// ProtocGenGoGRPC returns a singleton for protoc-gen-go-grpc.
func ProtocGenGoGRPC() *bintool.BinTool {
	if protocGenGoGRPC == nil {
		_ = Verdump().Ensure()
		ver := MustVersionLoadCache().GetVersion(ProtocGenGoGRPCVersion)
		loga.PrintInfo("Protocol Buffer Golang gRPC Version: %s", ver)
		protocGenGoGRPC = bintool.Must(bintool.NewGo(
			"google.golang.org/grpc/cmd/protoc-gen-go-grpc",
			ver,
			bintool.WithFolder(paths.MustGetProtobufPath()),
			bintool.WithVersionCmd(paths.MustGetGoBin("verdump")+" mod {{.FullCmd}}"),
		))
	}

	return protocGenGoGRPC
}

//nolint:gochecknoglobals // ignore globals
var protocGenGoTwirp *bintool.BinTool

// ProtocGenGoTwirp returns a singleton for protoc-gen-twirp.
func ProtocGenGoTwirp() *bintool.BinTool {
	if protocGenGoTwirp == nil {
		// ver := GetEnv("PROTOCGENGOTWIRP_VERSION", protocGenGoTwirpVersion)
		ver := MustVersionLoadCache().GetVersion(ProtocGenGoTwirpVersion)
		loga.PrintInfo("Protocol Buffer Golang Twirp Version: %s", ver)
		protocGenGoTwirp = bintool.Must(bintool.NewGo(
			"github.com/twitchtv/twirp/protoc-gen-twirp",
			ver,
			bintool.WithFolder(paths.MustGetProtobufPath()),
		))
	}

	return protocGenGoTwirp
}

//nolint:gochecknoglobals // ignore globals
var protocGenGoConnect *bintool.BinTool

// ProtocGenGoConnect returns a singleton for protoc-gen-connect-go.
func ProtocGenGoConnect() *bintool.BinTool {
	if protocGenGoConnect == nil {
		// ver := GetEnv("PROTOCGENGOCONNECT_VERSION", protocGenGoConnectVersion)
		ver := MustVersionLoadCache().GetVersion(ProtocGenGoConnectVersion)
		loga.PrintInfo("Protocol Buffer Golang ConnectRPC Version: %s", ver)
		protocGenGoConnect = bintool.Must(bintool.NewGo(
			"connectrpc.com/connect/cmd/protoc-gen-connect-go",
			ver,
			bintool.WithFolder(paths.MustGetProtobufPath()),
			bintool.WithVersionCmd(paths.MustGetGoBin("verdump")+" mod {{.FullCmd}}"),
		))
	}

	return protocGenGoConnect
}
