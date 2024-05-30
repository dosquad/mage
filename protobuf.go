package mage

import (
	"context"
	"os"
	"strings"

	"github.com/dosquad/mage/helper"
	"github.com/magefile/mage/mg"
	"github.com/princjef/mageutil/bintool"
)

func InstallProtoc(_ context.Context) error {
	return helper.BinProtoc().Ensure()
}

func InstallProtocGenGo(_ context.Context) error {
	return helper.BinProtocGenGo().Ensure()
}

func InstallProtocGenGoGRPC(_ context.Context) error {
	return helper.BinProtocGenGoGRPC().Ensure()
}

func InstallProtocGenGoTwirp(_ context.Context) error {
	return helper.BinProtocGenGoTwirp().Ensure()
}

func Protobuf(ctx context.Context) {
	mg.CtxDeps(ctx, InstallProtoc)
	mg.CtxDeps(ctx, InstallProtocGenGo)
	mg.CtxDeps(ctx, InstallProtocGenGoGRPC)
	mg.CtxDeps(ctx, ProtobufGenGo)
	mg.CtxDeps(ctx, ProtobufGenGoGRPC)
}

func ProtobufWithTwirp(ctx context.Context) {
	mg.CtxDeps(ctx, InstallProtoc)
	mg.CtxDeps(ctx, InstallProtocGenGo)
	mg.CtxDeps(ctx, InstallProtocGenGoGRPC)
	mg.CtxDeps(ctx, InstallProtocGenGoTwirp)
	mg.CtxDeps(ctx, ProtobufGenGo)
	mg.CtxDeps(ctx, ProtobufGenGoGRPC)
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
		if err := runProtoCommand(helper.BinProtoc(),
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
