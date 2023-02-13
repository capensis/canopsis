package js_test

import (
	"context"
	"errors"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/js"
	"github.com/dop251/goja"
	"github.com/kylelemons/godebug/pretty"
)

func TestCompile_GivenFunc_ShouldCallIt(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	code := `
        function test(a, b) {
          return a + b;
        }
	`
	e, err := js.Compile(t.Name(), code, "test")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	res, err := e.Execute(ctx, 1, 2)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	expected := 3
	if diff := pretty.Compare(res, expected); diff != "" {
		t.Fatalf("expected %v but got %v", expected, res)
	}
}

func TestCompile_GivenInvalidJs_ShouldReturnError(t *testing.T) {
	code := `
        function test(a, b) {
          return a + b;
	`
	_, err := js.Compile(t.Name(), code, "test")

	syntaxErr := &goja.CompilerSyntaxError{}
	if err == nil || !errors.As(err, &syntaxErr) {
		t.Fatalf("unexpected error %v", err)
	}
}

func TestCompile_GivenInvalidFuncArgs_ShouldReturnError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	code := `
        function test(a, b) {
          return a.indexOf(b);
        }
	`
	e, err := js.Compile(t.Name(), code, "test")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	_, err = e.Execute(ctx)
	exception := &goja.Exception{}
	if err == nil || !errors.As(err, &exception) {
		t.Fatalf("unexpected error %v", err)
	}
}
