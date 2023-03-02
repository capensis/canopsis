package js

import (
	"context"
	"fmt"

	"github.com/dop251/goja"
)

type Executor interface {
	Execute(ctx context.Context, args ...any) (any, error)
}

func Compile(name, src, funcName string) (Executor, error) {
	prg, err := goja.Compile(name, src, true)
	if err != nil {
		return nil, fmt.Errorf("cannot compile js: %w", err)
	}

	vm := goja.New()
	_, err = vm.RunProgram(prg)
	if err != nil {
		return nil, fmt.Errorf("cannot execute js: %w", err)
	}

	callable, ok := goja.AssertFunction(vm.Get(funcName))
	if !ok {
		return nil, fmt.Errorf("js function %q not found", funcName)
	}

	return &executor{
		vm:       vm,
		callable: callable,
	}, nil
}

type executor struct {
	vm       *goja.Runtime
	callable goja.Callable
}

func (e *executor) Execute(ctx context.Context, args ...any) (any, error) {
	transformedArgs := make([]goja.Value, len(args))
	for i, arg := range args {
		transformedArgs[i] = e.vm.ToValue(arg)
	}

	resCh := make(chan any, 1)
	go func() {
		defer close(resCh)
		r, err := e.callable(goja.Undefined(), transformedArgs...)
		if err != nil {
			resCh <- err
			return
		}

		resCh <- r.Export()
	}()

	select {
	case <-ctx.Done():
		return nil, nil
	case res := <-resCh:
		if err, ok := res.(error); ok {
			return nil, fmt.Errorf("cannot execute js function: %w", err)
		}

		return res, nil
	}
}
