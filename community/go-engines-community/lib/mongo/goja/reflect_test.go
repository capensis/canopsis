package goja

import (
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/dop251/goja"
	"github.com/kylelemons/godebug/pretty"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTransformOptions(t *testing.T) {
	for i, data := range getTestTransformOptionsDataSet() {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			prg, err := goja.Compile("test", data.JSCode, true)
			if err != nil {
				t.Fatalf("cannot compile js: %v", err)
			}

			vm := goja.New()
			err = vm.GlobalObject().Set("test", func(v goja.Value) {
				options := testOptions{}
				err := transformOptions(vm, v, &options)
				if data.ExpectedErr == "" && err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if data.ExpectedErr != "" && (err == nil || err.Error() != data.ExpectedErr) {
					t.Fatalf("expected error %q but got %q", data.ExpectedErr, err)
				}

				if diff := pretty.Compare(data.Expected, options); diff != "" {
					t.Fatalf("transformed value mismatch (-want +got):\n%s", diff)
				}
			})
			if err != nil {
				t.Fatalf("cannot set js func: %v", err)
			}

			_, err = vm.RunProgram(prg)
			if err != nil {
				t.Fatalf("cannot execute js: %v", err)
			}
		})
	}
}

func TestTransformValue(t *testing.T) {
	for i, data := range getTestTransformValueDataSet() {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			prg, err := goja.Compile("test", data.JSCode, true)
			if err != nil {
				t.Fatalf("cannot compile js: %v", err)
			}

			vm := goja.New()
			err = vm.GlobalObject().Set("test", func(v goja.Value) {
				result, err := transformValue(vm, v)
				if data.ExpectedErr == "" && err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if data.ExpectedErr != "" && (err == nil || err.Error() != data.ExpectedErr) {
					t.Fatalf("expected error %q but got %q", data.ExpectedErr, err)
				}

				if !reflect.DeepEqual(data.Expected, result) {
					if diff := pretty.Compare(data.Expected, result); diff != "" {
						t.Fatalf("transformed value mismatch (-want +got):\n%s", diff)
					} else {
						t.Fatalf("transformed value mismatch\nexpected:\n%[1]T %+[1]v\ngot:\n%[2]T %+[2]v", data.Expected, result)
					}
				}
			})
			if err != nil {
				t.Fatalf("cannot set js func: %v", err)
			}

			_, err = vm.RunProgram(prg)
			if err != nil {
				t.Fatalf("cannot execute js: %v", err)
			}
		})
	}
}

func getTestTransformOptionsDataSet() []struct {
	JSCode      string
	Expected    testOptions
	ExpectedErr string
} {
	timeInLocalTZ, err := time.ParseInLocation(time.DateTime, "2023-01-01 10:00:00", time.Local)
	if err != nil {
		panic(err)
	}

	return []struct {
		JSCode      string
		Expected    testOptions
		ExpectedErr string
	}{
		{
			JSCode:   `test({string: "test"})`,
			Expected: testOptions{str: "test"},
		},
		{
			JSCode:   `test({int64: 1})`,
			Expected: testOptions{i64: 1},
		},
		{
			JSCode:      `test({int32: 1})`,
			ExpectedErr: "set method for option \"int32\" accepts int32 instead of int64",
		},
		{
			JSCode:      `test({int64: 1.1})`,
			ExpectedErr: "set method for option \"int64\" accepts int64 instead of float64",
		},
		{
			JSCode:   `test({float64: 1.1})`,
			Expected: testOptions{f64: 1.1},
		},
		{
			JSCode:      `test({float32: 1.1})`,
			ExpectedErr: "set method for option \"float32\" accepts float32 instead of float64",
		},
		{
			JSCode:      `test({float64: 1.0})`,
			ExpectedErr: "set method for option \"float64\" accepts float64 instead of int64",
		},
		{
			JSCode:   `test({bool: true})`,
			Expected: testOptions{b: true},
		},
		{
			JSCode:   `test({time: new Date('2023-01-01T10:00:00')})`,
			Expected: testOptions{t: timeInLocalTZ},
		},
		{
			JSCode:      `test({duration: "30s"})`,
			ExpectedErr: "set method for option \"duration\" accepts time.Duration instead of string",
		},
		{
			JSCode:      `test({unknown: 1})`,
			ExpectedErr: "unknown option \"unknown\" for *goja.testOptions",
		},
		{
			JSCode:      `test({multiple: 1})`,
			ExpectedErr: "set method for option \"multiple\" has 2 arguments instead of one",
		},
	}
}

func getTestTransformValueDataSet() []struct {
	JSCode      string
	Expected    any
	ExpectedErr string
} {
	timeInLocalTZ, err := time.ParseInLocation(time.DateTime, "2023-01-01 10:00:00", time.Local)
	if err != nil {
		panic(err)
	}

	return []struct {
		JSCode      string
		Expected    any
		ExpectedErr string
	}{
		{
			JSCode:   `test(1)`,
			Expected: int64(1),
		},
		{
			JSCode:   `test(1.0)`,
			Expected: int64(1.0),
		},
		{
			JSCode:   `test(1.1)`,
			Expected: float64(1.1),
		},
		{
			JSCode:   `test("smth")`,
			Expected: "smth",
		},
		{
			JSCode:   `test(true)`,
			Expected: true,
		},
		{
			JSCode:   `test(null)`,
			Expected: nil,
		},
		{
			JSCode:   `test({foo: "bar"})`,
			Expected: bson.D{{Key: "foo", Value: "bar"}},
		},
		{
			JSCode:   `test({foo: {bar: "baz"}})`,
			Expected: bson.D{{Key: "foo", Value: bson.D{{Key: "bar", Value: "baz"}}}},
		},
		{
			JSCode:   `test({foo: "bar", baz: "aux"})`,
			Expected: bson.D{{Key: "foo", Value: "bar"}, {Key: "baz", Value: "aux"}},
		},
		{
			JSCode:   `test([1, 2])`,
			Expected: bson.A{int64(1), int64(2)},
		},
		{
			JSCode:   `test(new Array(1, 2))`,
			Expected: bson.A{int64(1), int64(2)},
		},
		{
			JSCode:   `test([{foo: "bar"}, {baz: "aux"}])`,
			Expected: bson.A{bson.D{{Key: "foo", Value: "bar"}}, bson.D{{Key: "baz", Value: "aux"}}},
		},
		{
			JSCode:   `test(new Map([["foo", "bar"], ["baz", "auz"]]))`,
			Expected: bson.A{bson.A{"foo", "bar"}, bson.A{"baz", "auz"}},
		},
		{
			JSCode:   `test(new Set([1, 2, 3]))`,
			Expected: bson.A{int64(1), int64(2), int64(3)},
		},
		{
			JSCode:   `test(new Date('2023-01-01T10:00:00'))`,
			Expected: timeInLocalTZ,
		},
		{
			JSCode:   `test(new Date('2023-01-01T16:00:00+07:00'))`,
			Expected: timeInLocalTZ,
		},
		{
			JSCode: `test(/^foo/i)`,
			Expected: primitive.Regex{
				Pattern: "^foo",
				Options: "i",
			},
		},
		{
			JSCode:      `test(new Promise(function(resolve, reject) {}))`,
			Expected:    nil,
			ExpectedErr: "unsupported value type: *goja.Promise",
		},
	}
}

type testOptions struct {
	str string
	i64 int64
	i32 int32
	f64 float64
	f32 float32
	b   bool
	t   time.Time
	d   time.Duration
}

func (o *testOptions) SetString(v string) {
	o.str = v
}

func (o *testOptions) SetInt64(v int64) {
	o.i64 = v
}

func (o *testOptions) SetInt32(v int32) {
	o.i32 = v
}

func (o *testOptions) SetFloat64(v float64) {
	o.f64 = v
}

func (o *testOptions) SetFloat32(v float32) {
	o.f32 = v
}

func (o *testOptions) SetBool(v bool) {
	o.b = v
}

func (o *testOptions) SetTime(v time.Time) {
	o.t = v
}

func (o *testOptions) SetDuration(v time.Duration) {
	o.d = v
}

func (o *testOptions) SetMultiple(_, _ int64) {
}
