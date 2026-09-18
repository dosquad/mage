package builder

// validPlatforms contains the official list of supported GOOS/GOARCH pairs.
//
//nolint:gochecknoglobals,goconst // constant would reduce readability in this case.
var validPlatforms = map[string]map[string]bool{
	"aix":       {"ppc64": true},
	"android":   {"386": true, "amd64": true, "arm": true, "arm64": true},
	"darwin":    {"amd64": true, "arm64": true},
	"dragonfly": {"amd64": true},
	"freebsd":   {"386": true, "amd64": true, "arm": true, "arm64": true, "riscv64": true},
	"illumos":   {"amd64": true},
	"js":        {"wasm": true},
	"wasip1":    {"wasm": true},
	"linux": {
		"386":      true,
		"amd64":    true,
		"arm":      true,
		"arm64":    true,
		"loong64":  true,
		"mips":     true,
		"mipsle":   true,
		"mips64":   true,
		"mips64le": true,
		"ppc64":    true,
		"ppc64le":  true,
		"riscv64":  true,
		"s390x":    true,
	},
	"netbsd":  {"386": true, "amd64": true, "arm": true, "arm64": true},
	"openbsd": {"386": true, "amd64": true, "arm": true, "arm64": true, "riscv64": true},
	"plan9":   {"386": true, "amd64": true, "arm": true},
	"solaris": {"amd64": true},
	"windows": {"386": true, "amd64": true, "arm": true, "arm64": true},
}

// IsValidTarget verifies if the GOOS and GOARCH combination is valid.
func IsValidTarget(goos, goarch string) bool {
	archs, osExists := validPlatforms[goos]
	if !osExists {
		return false
	}
	return archs[goarch]
}
