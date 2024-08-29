package dyndep

import "context"

//nolint:gochecknoglobals // singleton for dynamic dependencies.
var dMap *dynDep

func initMap() {
	if dMap == nil {
		dMap = newDynDep()
	}
}

func Add(target Type, f interface{}) {
	initMap()

	dMap.Add(target, f)
}

func CtxDeps(ctx context.Context, target Type) {
	initMap()

	dMap.CtxDeps(ctx, target)
}
