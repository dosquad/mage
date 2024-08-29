package dyndep

import (
	"context"
	"sync"

	"github.com/magefile/mage/mg"
)

type dynDep struct {
	lock sync.Mutex
	deps map[Type][]interface{}
}

func newDynDep() *dynDep {
	return &dynDep{
		deps: make(map[Type][]interface{}),
	}
}

func (d *dynDep) Add(target Type, f interface{}) {
	d.lock.Lock()
	defer d.lock.Unlock()

	if len(d.deps[target]) == 0 {
		d.deps[target] = []interface{}{}
	}

	d.deps[target] = append(d.deps[target], f)
}

func (d *dynDep) CtxDeps(ctx context.Context, target Type) {
	d.lock.Lock()
	defer d.lock.Unlock()

	for _, dep := range d.deps[target] {
		mg.CtxDeps(ctx, dep)
	}

	d.deps[target] = []interface{}{}
}
