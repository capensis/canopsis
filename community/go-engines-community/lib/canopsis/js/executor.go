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
		return nil, err
	}

	vm := goja.New()
	_, err = vm.RunProgram(prg)
	if err != nil {
		return nil, err
	}

	callable, ok := goja.AssertFunction(vm.Get(funcName))
	if !ok {
		return nil, fmt.Errorf("%s not found", funcName)
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

	stop := make(chan struct{})
	defer close(stop)

	go func() {
		select {
		case <-stop:
		case <-ctx.Done():
			e.vm.Interrupt(ctx.Err())
		}
	}()

	r, err := e.callable(goja.Undefined(), transformedArgs...)
	if err != nil {
		return nil, err
	}

	return r.Export(), nil
}
