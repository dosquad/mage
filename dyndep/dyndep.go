package dyndep

import (
	"context"
	"sync"

	"github.com/dosquad/mage/loga"
	"github.com/magefile/mage/mg"
)

type dynDep struct {
	lock sync.Mutex
	deps map[Type][]any
}

func newDynDep() *dynDep {
	return &dynDep{
		deps: make(map[Type][]any),
	}
}

func (d *dynDep) Add(target Type, f any) {
	loga.PrintDebugf("dynDep.Add(%s): dep=%T", target, f)
	d.lock.Lock()
	defer d.lock.Unlock()

	if len(d.deps[target]) == 0 {
		d.deps[target] = []any{}
	}

	d.deps[target] = append(d.deps[target], f)
}

func (d *dynDep) CtxDeps(ctx context.Context, target Type) {
	d.lock.Lock()
	defer d.lock.Unlock()

	for _, dep := range d.deps[target] {
		mg.CtxDeps(ctx, dep)
	}

	d.deps[target] = []any{}
}

func (d *dynDep) Get(target Type) []any {
	d.lock.Lock()
	defer d.lock.Unlock()

	out := make([]any, 0, len(d.deps[target]))
	for _, dep := range d.deps[target] {
		loga.PrintDebugf("dynDep.Get(%s): dep=%T", target, dep)
		out = append(out, dep)
	}

	return out
}
