package dyndep

import (
	"context"

	"github.com/dosquad/mage/loga"
)

type Type string

const (
	ArchiveSideload Type = "archive-sideload"
	Build           Type = "build"
	Clean           Type = "clean"
	Cfssl           Type = "cfssl"
	Docker          Type = "docker"
	Golang          Type = "golang"
	Goreleaser      Type = "goreleaser"
	Install         Type = "install"
	Kubebuilder     Type = "kubebuilder"
	Lint            Type = "lint"
	Mirrord         Type = "mirrord"
	Mod             Type = "mod"
	Protobuf        Type = "protobuf"
	Run             Type = "run"
	Test            Type = "test"
	Update          Type = "update"
	Wire            Type = "wire"
)

type ArchiveSideloadFunc func(context.Context) map[string]string

func GetArchiveSideloadDeps(ctx context.Context) map[string]string {
	initMap()

	// Retrieve the dependencies for the ArchiveSideload type
	deps := dMap.Get(ArchiveSideload)

	loga.PrintDebugf("deps: %+v", deps)

	// Combine all dependencies into a single map
	result := make(map[string]string)
	for _, dep := range deps {
		loga.PrintDebugf("dep: %+v", dep)
		if fn, ok := dep.(ArchiveSideloadFunc); ok {
			loga.PrintDebugf("dep:ok: %t", ok)
			for k, v := range fn(ctx) {
				loga.PrintDebugf("Adding to result: key=%s, value=%s", k, v)
				result[k] = v
			}
		}
	}

	return result
}
