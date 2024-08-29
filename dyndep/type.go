package dyndep

type Type string

const (
	Build       Type = "build"
	Clean       Type = "clean"
	Cfssl       Type = "cfssl"
	Docker      Type = "docker"
	Golang      Type = "golang"
	Goreleaser  Type = "goreleaser"
	Install     Type = "install"
	Kubebuilder Type = "kubebuilder"
	Lint        Type = "lint"
	Mirrord     Type = "mirrord"
	Mod         Type = "mod"
	Protobuf    Type = "protobuf"
	Run         Type = "run"
	Test        Type = "test"
	Update      Type = "update"
	Wire        Type = "wire"
)
