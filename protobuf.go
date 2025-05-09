package mage

import (
	"context"
	"os"
	"strings"

	"github.com/dosquad/mage/dyndep"
	"github.com/dosquad/mage/helper/bins"
	"github.com/dosquad/mage/helper/build"
	"github.com/dosquad/mage/helper/must"
	"github.com/dosquad/mage/helper/paths"
	"github.com/dosquad/mage/helper/pbuf"
	"github.com/magefile/mage/mg"
	"github.com/princjef/mageutil/bintool"
)

// Protobuf namespace is defined to group Protocol buffer functions.
type Protobuf mg.Namespace

// installProtoc install protoc command.
func (Protobuf) installProtoc(_ context.Context) error {
	return bins.Protoc().Ensure()
}

// installProtocGenGo install protoc-gen-go command.
func (Protobuf) installProtocGenGo(_ context.Context) error {
	return bins.ProtocGenGo().Ensure()
}

// installProtocGenGoGRPC install protoc-gen-go-grpc command.
func (Protobuf) installProtocGenGoGRPC(_ context.Context) error {
	return bins.ProtocGenGoGRPC().Ensure()
}

// installProtocGenGoTwirp install protoc-gen-go-twirp command.
func (Protobuf) installProtocGenGoTwirp(_ context.Context) error {
	return bins.ProtocGenGoTwirp().Ensure()
}

func (Protobuf) installProtocGenGoConnect(_ context.Context) error {
	return bins.ProtocGenGoConnect().Ensure()
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

// GenerateWithConnect install and generate golang Protocol Buffer files (including ConnectRPC).
func (Protobuf) GenerateWithConnect(ctx context.Context) {
	dyndep.CtxDeps(ctx, dyndep.Protobuf)

	mg.CtxDeps(ctx, Protobuf.installProtoc)
	mg.CtxDeps(ctx, Protobuf.installProtocGenGo)
	mg.CtxDeps(ctx, Protobuf.installProtocGenGoGRPC)
	mg.CtxDeps(ctx, Protobuf.installProtocGenGoConnect)
	mg.CtxDeps(ctx, Protobuf.GenGo)
	mg.CtxDeps(ctx, Protobuf.GenGoGRPC)
	mg.CtxDeps(ctx, Protobuf.GenGoConnect)
}

func runProtoCommand(cmd *bintool.BinTool, args []string) error {
	origPath := os.Getenv("PATH")
	defer func() { _ = os.Setenv("PATH", origPath) }()

	if err := os.Setenv("PATH", paths.MustGetProtobufPath()+":"+origPath); err != nil {
		return err
	}

	// loga.PrintInfo("runProtoCommand: PATH=%s", os.Getenv("PATH"))

	return cmd.Command(strings.Join(args, " ")).Run()
}

// GenGo run protoc-gen-go to generate code.
func (Protobuf) GenGo(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Protobuf)

	mg.CtxDeps(ctx, Protobuf.installProtoc)
	mg.CtxDeps(ctx, Protobuf.installProtocGenGo)

	return protobufGen(ctx, []string{
		"--proto_path=" + paths.MustGetArtifactPath("protobuf", "include"),
		"--go_opt=module=" + must.Must[string](build.GetModuleName()),
		"--go_out=.",
	})
}

// GenGoGRPC run protoc-gen-go-grpc to generate code.
func (Protobuf) GenGoGRPC(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Protobuf)

	mg.CtxDeps(ctx, Protobuf.installProtoc)
	mg.CtxDeps(ctx, Protobuf.installProtocGenGoGRPC)

	return protobufGen(ctx, []string{
		"--proto_path=" + paths.MustGetArtifactPath("protobuf", "include"),
		"--go-grpc_opt=module=" + must.Must[string](build.GetModuleName()),
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
		"--proto_path=" + paths.MustGetArtifactPath("protobuf", "include"),
		"--go_opt=module=" + must.Must[string](build.GetModuleName()),
		"--go_out=.",
		"--twirp_opt=module=" + must.Must[string](build.GetModuleName()),
		"--twirp_out=.",
	})
}

// GenGoConnect run protoc-gen-connect-go to generate code.
func (Protobuf) GenGoConnect(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Protobuf)

	mg.CtxDeps(ctx, Protobuf.installProtoc)
	mg.CtxDeps(ctx, Protobuf.installProtocGenGo)
	mg.CtxDeps(ctx, Protobuf.installProtocGenGoConnect)

	return protobufGen(ctx, []string{
		"--proto_path=" + paths.MustGetArtifactPath("protobuf", "include"),
		"--go_opt=module=" + must.Must[string](build.GetModuleName()),
		"--go_out=.",
		"--connect-go_opt=module=" + must.Must[string](build.GetModuleName()),
		"--connect-go_out=.",
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
	protobufPaths := pbuf.ProtobufIncludePaths()

	for _, protoPathFunc := range pbuf.ProtobufTargets() {
		if err := runProtoCommand(bins.Protoc(),
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
