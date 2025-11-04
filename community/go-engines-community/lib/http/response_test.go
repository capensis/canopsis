package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"reflect"
	"strconv"
	"testing"

	"github.com/valyala/fastjson"
)

func TestReadResponse(t *testing.T) {
	longBuff := make([]byte, 2*buffChunk+buffChunk/2)
	for i := 0; i < len(longBuff); i++ {
		longBuff[i] = byte('a' + rand.Intn(26))
	}

	dataSet := []struct {
		Body        io.ReadCloser
		MaxSize     int64
		ExpectedRes []byte
		ExpectedErr error
	}{
		{
			Body:        io.NopCloser(bytes.NewReader(longBuff)),
			MaxSize:     int64(len(longBuff)),
			ExpectedRes: longBuff,
		},
		{
			Body:        io.NopCloser(bytes.NewReader(longBuff[:buffChunk/2])),
			MaxSize:     buffChunk,
			ExpectedRes: longBuff[:buffChunk/2],
		},
		{
			ExpectedErr: ErrResponseTooLong,
		},
		{
			Body:        io.NopCloser(bytes.NewReader(longBuff)),
			MaxSize:     int64(len(longBuff) / 4 * 3),
			ExpectedErr: ErrResponseTooLong,
		},
	}

	for i, data := range dataSet {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			res, err := ReadResponse(data.Body, data.MaxSize)
			if !errors.Is(err, data.ExpectedErr) {
				t.Errorf("expected err %v but got %v", data.ExpectedErr, err)
			}
			if string(res) != string(data.ExpectedRes) {
				t.Errorf("expected result\n%s\nbut got\n%s", data.ExpectedRes, res)
			}
		})
	}
}

func TestFlatten(t *testing.T) {
	dataSets := map[string]struct {
		input          any
		expectedOutput map[string]any
	}{
		"test simple map": {
			input: map[string]any{
				"a": 1,
				"b": "test",
				"c": false,
			},
			expectedOutput: map[string]any{
				"a": 1,
				"b": "test",
				"c": false,
			},
		},
		"test map with nested maps": {
			input: map[string]any{
				"a": 1,
				"b": map[string]any{
					"d": 1,
					"e": map[string]any{
						"f": "test",
					},
				},
				"c": 2,
			},
			expectedOutput: map[string]any{
				"a":     1,
				"b.d":   1,
				"b.e.f": "test",
				"c":     2,
			},
		},
		"test map with arrays": {
			input: map[string]any{
				"a": []any{
					map[string]any{
						"b": 1,
						"c": "test",
					},
					map[string]any{
						"b": 2,
						"c": "test 2",
					},
					map[string]any{
						"b": 3,
						"c": map[string]any{
							"d": 2,
							"e": "test",
						},
					},
				},
				"f": map[string]any{
					"g": 3,
					"h": "test",
				},
				"i": []any{
					map[string]any{
						"j": 1,
						"k": []any{
							map[string]any{
								"l": 10,
								"m": true,
							},
							map[string]any{
								"l": 20,
								"m": false,
							},
							map[string]any{
								"l": 30,
								"m": true,
							},
						},
					},
					map[string]any{
						"j": 2,
						"k": []any{
							map[string]any{
								"l": 30,
								"m": false,
							},
							map[string]any{
								"l": 20,
								"m": true,
							},
							map[string]any{
								"l": 10,
								"m": false,
							},
						},
					},
				},
				"n": []any{
					1,
					2,
					3,
				},
				"o": []any{
					1,
					"2",
					true,
					nil,
				},
				"p": []any{
					[]any{
						map[string]any{
							"q": 1,
							"r": []any{
								"1",
								"2",
								"3",
							},
							"s": 3,
						},
					},
				},
				"t": []any{},
				"u": []any{[]any{}},
			},
			expectedOutput: map[string]any{
				"a.0.b":     1,
				"a.0.c":     "test",
				"a.1.b":     2,
				"a.1.c":     "test 2",
				"a.2.b":     3,
				"a.2.c.d":   2,
				"a.2.c.e":   "test",
				"f.g":       3,
				"f.h":       "test",
				"i.0.j":     1,
				"i.1.j":     2,
				"i.0.k.0.l": 10,
				"i.0.k.0.m": true,
				"i.0.k.1.l": 20,
				"i.0.k.1.m": false,
				"i.0.k.2.l": 30,
				"i.0.k.2.m": true,
				"i.1.k.0.l": 30,
				"i.1.k.0.m": false,
				"i.1.k.1.l": 20,
				"i.1.k.1.m": true,
				"i.1.k.2.l": 10,
				"i.1.k.2.m": false,
				"n":         []any{1, 2, 3},
				"n.0":       1,
				"n.1":       2,
				"n.2":       3,
				"o":         []any{1, "2", true, nil},
				"o.0":       1,
				"o.1":       "2",
				"o.2":       true,
				"o.3":       nil,
				"p.0.0.q":   1,
				"p.0.0.r":   []any{"1", "2", "3"},
				"p.0.0.r.0": "1",
				"p.0.0.r.1": "2",
				"p.0.0.r.2": "3",
				"p.0.0.s":   3,
				"t":         []any{},
				"u.0":       []any{},
			},
		},
		"test_input_array": {
			input: []any{
				1,
				"abc",
				map[string]any{
					"a": 1,
					"b": 2,
					"c": 3,
				},
				[]any{
					4,
					5,
					"6",
				},
			},
			expectedOutput: map[string]any{
				"0":   1,
				"1":   "abc",
				"2.a": 1,
				"2.b": 2,
				"2.c": 3,
				"3.0": 4,
				"3.1": 5,
				"3.2": "6",
				"3":   []any{4, 5, "6"},
			},
		},
	}

	for test, dataSet := range dataSets {
		t.Run(test, func(t *testing.T) {
			b, _ := json.Marshal(dataSet.input)
			v, _ := fastjson.ParseBytes(b)
			result := flatten(v, "")
			if !reflect.DeepEqual(result, dataSet.expectedOutput) {
				t.Errorf("expected media %v but got %v", dataSet.expectedOutput, result)
			}
		})
	}
}
