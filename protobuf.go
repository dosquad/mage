package mage

import (
	"context"
	"os"
	"runtime"
	"strings"

	"github.com/dosquad/mage/helper"
	"github.com/magefile/mage/mg"
	"github.com/princjef/mageutil/bintool"
)

//nolint:gochecknoglobals // ignore globals
var protoc *bintool.BinTool

func installProtoc() *bintool.BinTool {
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

		var protocVer string
		{
			var err error
			protocVer, err = helper.HTTPGetLatestGitHubVersion("protocolbuffers/protobuf")
			if err != nil {
				helper.PrintWarning("Protocol Buffer Error: %s", err)
				protocVer = "latest"
			}
		}

		helper.PrintInfo("Protocol Buffer Version: %s", protocVer)
		helper.PanicIfError(helper.ExtractArchive(
			"https://github.com/protocolbuffers/protobuf/releases/download/v"+protocVer+"/"+
				"protoc-"+protocVer+"-"+goOperatingSystem+"-"+goArch+".zip",
			helper.MustGetWD("artifacts", "protobuf"),
		), "Extract Archive")

		protoc = bintool.Must(bintool.New(
			"protoc{{.BinExt}}",
			protocVer,
			"https://github.com/protocolbuffers/protobuf/releases/download/v"+protocVer+"/"+
				"protoc-"+protocVer+"-"+goOperatingSystem+"-"+goArch+".zip",
			bintool.WithFolder(helper.MustGetWD("artifacts", "protobuf", "bin")),
		))
	}

	return protoc
}
func InstallProtoc(_ context.Context) error {
	return installProtoc().Ensure()
}

//nolint:gochecknoglobals // ignore globals
var protocGenGo *bintool.BinTool

func installProtocGenGo() *bintool.BinTool {
	if protocGenGo == nil {
		helper.PrintInfo("Protocol Buffer Golang Version: %s", helper.GetProtobufVersion())
		protocGenGo = bintool.Must(bintool.NewGo(
			"google.golang.org/protobuf/cmd/protoc-gen-go",
			helper.GetProtobufVersion(),
			bintool.WithFolder(helper.MustGetProtobufPath()),
		))
	}

	return protocGenGo
}
func InstallProtocGenGo(_ context.Context) error {
	return installProtocGenGo().Ensure()
}

//nolint:gochecknoglobals // ignore globals
var protocGenGoGRPC *bintool.BinTool

func installProtocGenGoGRPC() *bintool.BinTool {
	if protocGenGoGRPC == nil {
		protocGenGoGRPC = bintool.Must(bintool.NewGo(
			"google.golang.org/grpc/cmd/protoc-gen-go-grpc",
			protocGenGoGRPCVersion,
			bintool.WithFolder(helper.MustGetProtobufPath()),
		))
	}

	return protocGenGoGRPC
}
func InstallProtocGenGoGRPC(_ context.Context) error {
	return installProtocGenGoGRPC().Ensure()
}

//nolint:gochecknoglobals // ignore globals
var protocGenGoTwirp *bintool.BinTool

func installProtocGenGoTwirp() *bintool.BinTool {
	if protocGenGoTwirp == nil {
		protocGenGoTwirp = bintool.Must(bintool.NewGo(
			"github.com/twitchtv/twirp/protoc-gen-twirp",
			protocGenGoTwirpVersion,
			bintool.WithFolder(helper.MustGetProtobufPath()),
		))
	}

	return protocGenGoTwirp
}
func InstallProtocGenGoTwirp(_ context.Context) error {
	return installProtocGenGoTwirp().Ensure()
}

func Protobuf(ctx context.Context) {
	mg.CtxDeps(ctx, InstallProtoc)
	mg.CtxDeps(ctx, InstallProtocGenGo)
	mg.CtxDeps(ctx, InstallProtocGenGoGRPC)
	mg.CtxDeps(ctx, ProtobufGenGoGRPC)
	mg.CtxDeps(ctx, ProtobufGenGo)
}

func ProtobufWithTwirp(ctx context.Context) {
	mg.CtxDeps(ctx, InstallProtoc)
	mg.CtxDeps(ctx, InstallProtocGenGo)
	mg.CtxDeps(ctx, InstallProtocGenGoGRPC)
	mg.CtxDeps(ctx, InstallProtocGenGoTwirp)
	mg.CtxDeps(ctx, ProtobufGenGoGRPC)
	mg.CtxDeps(ctx, ProtobufGenGo)
	mg.CtxDeps(ctx, ProtobufGenGoTwirp)
}

func runProtoCommand(cmd *bintool.BinTool, args []string) error {
	origPath := os.Getenv("PATH")
	defer func() { os.Setenv("PATH", origPath) }()

	os.Setenv("PATH", helper.MustGetProtobufPath()+":"+origPath)

	return cmd.Command(strings.Join(args, " ")).Run()
}

func ProtobufGenGo(ctx context.Context) error {
	mg.CtxDeps(ctx, InstallProtoc)
	mg.CtxDeps(ctx, InstallProtocGenGo)

	return protobufGen(ctx, []string{
		"--proto_path=" + helper.MustGetWD("artifacts", "protobuf", "include"),
		"--go_opt=module=" + helper.MustModuleName(),
		"--go_out=.",
	})
}

func ProtobufGenGoGRPC(ctx context.Context) error {
	mg.CtxDeps(ctx, InstallProtoc)
	mg.CtxDeps(ctx, InstallProtocGenGoGRPC)

	return protobufGen(ctx, []string{
		"--proto_path=" + helper.MustGetWD("artifacts", "protobuf", "include"),
		"--go-grpc_opt=module=" + helper.MustModuleName(),
		"--go-grpc_out=.",
		"--go-grpc_opt=require_unimplemented_servers=false",
	})
}

func ProtobufGenGoTwirp(ctx context.Context) error {
	mg.CtxDeps(ctx, InstallProtoc)
	mg.CtxDeps(ctx, InstallProtocGenGo)
	mg.CtxDeps(ctx, InstallProtocGenGoTwirp)

	return protobufGen(ctx, []string{
		"--proto_path=" + helper.MustGetWD("artifacts", "protobuf", "include"),
		"--go_opt=module=" + helper.MustModuleName(),
		"--go_out=.",
		"--twirp_opt=module=" + helper.MustModuleName(),
		"--twirp_out=.",
	})
}

func protobufGen(_ context.Context, coreArgs []string) error {
	// mg.CtxDeps(ctx, InstallProtoc)
	// mg.CtxDeps(ctx, InstallProtocGenGoGRPC)

	// var moduleName string
	// {
	// 	var err error
	// 	moduleName, err = helper.GetModuleName()
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// coreArgs := []string{
	// 	"--proto_path=" + helper.MustGetWD("artifacts", "protobuf", "include"),
	// 	"--go-grpc_opt=module=" + moduleName,
	// 	"--go-grpc_out=.",
	// 	"--go-grpc_opt=require_unimplemented_servers=false",
	// }
	protobufPaths := helper.ProtobufIncludePaths()

	for _, protoPathFunc := range helper.ProtobufTargets() {
		if err := runProtoCommand(protoc,
			append(
				append(coreArgs, protobufPaths...),
				" "+strings.Join(protoPathFunc(), " "),
			),
		); err != nil {
			return err
		}
	}

	return nil
}
