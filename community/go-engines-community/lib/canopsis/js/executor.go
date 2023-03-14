package js

import (
	"context"
	"fmt"

	"github.com/dop251/goja"
)

type Executor interface {
	ExecuteFunc(ctx context.Context, funcName string, args ...any) (any, error)
}

func Compile(name, src string) (Executor, error) {
	prg, err := goja.Compile(name, src, true)
	if err != nil {
		return nil, fmt.Errorf("cannot compile js: %w", err)
	}

	return &executor{prg: prg}, nil
}

type executor struct {
	prg *goja.Program
}

func (e *executor) ExecuteFunc(ctx context.Context, funcName string, args ...any) (any, error) {
	vm := goja.New()
	_, err := vm.RunProgram(e.prg)
	if err != nil {
		return nil, fmt.Errorf("cannot execute js: %w", err)
	}

	callable, ok := goja.AssertFunction(vm.Get(funcName))
	if !ok {
		return nil, fmt.Errorf("js function %q not found", funcName)
	}

	transformedArgs := make([]goja.Value, len(args))
	for i, arg := range args {
		transformedArgs[i] = vm.ToValue(arg)
	}

	resCh := make(chan any, 1)
	go func() {
		defer close(resCh)
		r, err := callable(goja.Undefined(), transformedArgs...)
		if err != nil {
			resCh <- err
			return
		}

		resCh <- r.Export()
	}()

	select {
	case <-ctx.Done():
		vm.Interrupt(ctx.Err())
		return nil, ctx.Err()
	case res := <-resCh:
		if err, ok := res.(error); ok {
			return nil, fmt.Errorf("cannot execute js function: %w", err)
		}

		return res, nil
	}
}
