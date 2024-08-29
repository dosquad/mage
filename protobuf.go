package mage

import (
	"context"
	"os"
	"strings"

	"github.com/dosquad/mage/dyndep"
	"github.com/dosquad/mage/helper"
	"github.com/magefile/mage/mg"
	"github.com/princjef/mageutil/bintool"
)

// Protobuf namespace is defined to group Protocol buffer functions.
type Protobuf mg.Namespace

// installProtoc install protoc command.
func (Protobuf) installProtoc(_ context.Context) error {
	return helper.BinProtoc().Ensure()
}

// installProtocGenGo install protoc-gen-go command.
func (Protobuf) installProtocGenGo(_ context.Context) error {
	return helper.BinProtocGenGo().Ensure()
}

// installProtocGenGoGRPC install protoc-gen-go-grpc command.
func (Protobuf) installProtocGenGoGRPC(_ context.Context) error {
	return helper.BinProtocGenGoGRPC().Ensure()
}

// installProtocGenGoTwirp install protoc-gen-go-twirp command.
func (Protobuf) installProtocGenGoTwirp(_ context.Context) error {
	return helper.BinProtocGenGoTwirp().Ensure()
}

// Generate install and generate golang Protocol Buffer files.
func (Protobuf) Generate(ctx context.Context) {
	dyndep.CtxDeps(ctx, dyndep.Protobuf)

	mg.CtxDeps(ctx, Protobuf.installProtoc)
	mg.CtxDeps(ctx, Protobuf.installProtocGenGo)
	mg.CtxDeps(ctx, Protobuf.installProtocGenGoGRPC)
	mg.CtxDeps(ctx, Protobuf.GenGo)
	mg.CtxDeps(ctx, Protobuf.GenGoGRPC)
}

// GenerateWithTwirp install and generate golang Protocol Buffer files (including Twirp).
func (Protobuf) GenerateWithTwirp(ctx context.Context) {
	dyndep.CtxDeps(ctx, dyndep.Protobuf)

	mg.CtxDeps(ctx, Protobuf.installProtoc)
	mg.CtxDeps(ctx, Protobuf.installProtocGenGo)
	mg.CtxDeps(ctx, Protobuf.installProtocGenGoGRPC)
	mg.CtxDeps(ctx, Protobuf.installProtocGenGoTwirp)
	mg.CtxDeps(ctx, Protobuf.GenGo)
	mg.CtxDeps(ctx, Protobuf.GenGoGRPC)
	mg.CtxDeps(ctx, Protobuf.GenGoTwirp)
}

func runProtoCommand(cmd *bintool.BinTool, args []string) error {
	origPath := os.Getenv("PATH")
	defer func() { os.Setenv("PATH", origPath) }()

	os.Setenv("PATH", helper.MustGetProtobufPath()+":"+origPath)

	return cmd.Command(strings.Join(args, " ")).Run()
}

// ProtobufGenGo run protoc-gen-go to generate code.
func (Protobuf) GenGo(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Protobuf)

	mg.CtxDeps(ctx, Protobuf.installProtoc)
	mg.CtxDeps(ctx, Protobuf.installProtocGenGo)

	return protobufGen(ctx, []string{
		"--proto_path=" + helper.MustGetArtifactPath("protobuf", "include"),
		"--go_opt=module=" + helper.Must[string](helper.GetModuleName()),
		"--go_out=.",
	})
}

// GenGoGRPC run protoc-gen-go-grpc to generate code.
func (Protobuf) GenGoGRPC(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Protobuf)

	mg.CtxDeps(ctx, Protobuf.installProtoc)
	mg.CtxDeps(ctx, Protobuf.installProtocGenGoGRPC)

	return protobufGen(ctx, []string{
		"--proto_path=" + helper.MustGetArtifactPath("protobuf", "include"),
		"--go-grpc_opt=module=" + helper.Must[string](helper.GetModuleName()),
		"--go-grpc_out=.",
		"--go-grpc_opt=require_unimplemented_servers=false",
	})
}

// GenGoTwirp run protoc-gen-go-twirp to generate code.
func (Protobuf) GenGoTwirp(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Protobuf)

	mg.CtxDeps(ctx, Protobuf.installProtoc)
	mg.CtxDeps(ctx, Protobuf.installProtocGenGo)
	mg.CtxDeps(ctx, Protobuf.installProtocGenGoTwirp)

	return protobufGen(ctx, []string{
		"--proto_path=" + helper.MustGetArtifactPath("protobuf", "include"),
		"--go_opt=module=" + helper.Must[string](helper.GetModuleName()),
		"--go_out=.",
		"--twirp_opt=module=" + helper.Must[string](helper.GetModuleName()),
		"--twirp_out=.",
	})
}

func protobufGen(_ context.Context, coreArgs []string) error {
	// mg.CtxDeps(ctx, Protobuf.installProtoc)
	// mg.CtxDeps(ctx, Protobuf.installProtocGenGoGRPC)

	// var moduleName string
	// {
	// 	var err error
	// 	moduleName, err = helper.GetModuleName()
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// coreArgs := []string{
	// 	"--proto_path=" + helper.MustGetArtifactPath("protobuf", "include"),
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
