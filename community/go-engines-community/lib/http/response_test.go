package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net/http"
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
		Response    *http.Response
		MaxSize     int64
		ExpectedRes []byte
		ExpectedErr error
	}{
		{
			Response: &http.Response{
				Body: io.NopCloser(bytes.NewReader(longBuff)),
			},
			MaxSize:     int64(len(longBuff)),
			ExpectedRes: longBuff,
		},
		{
			Response: &http.Response{
				Body: io.NopCloser(bytes.NewReader(longBuff[:buffChunk/2])),
			},
			MaxSize:     buffChunk,
			ExpectedRes: longBuff[:buffChunk/2],
		},
		{
			ExpectedErr: ErrResponseTooLong,
		},
		{
			Response: &http.Response{
				Body: io.NopCloser(bytes.NewReader(longBuff)),
			},
			MaxSize:     int64(len(longBuff) / 4 * 3),
			ExpectedErr: ErrResponseTooLong,
		},
	}

	for i, data := range dataSet {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			res, err := ReadResponse(data.Response, data.MaxSize)
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
		input          interface{}
		expectedOutput map[string]interface{}
	}{
		"test simple map": {
			input: map[string]interface{}{
				"a": 1,
				"b": "test",
				"c": false,
			},
			expectedOutput: map[string]interface{}{
				"a": 1,
				"b": "test",
				"c": false,
			},
		},
		"test map with nested maps": {
			input: map[string]interface{}{
				"a": 1,
				"b": map[string]interface{}{
					"d": 1,
					"e": map[string]interface{}{
						"f": "test",
					},
				},
				"c": 2,
			},
			expectedOutput: map[string]interface{}{
				"a":     1,
				"b.d":   1,
				"b.e.f": "test",
				"c":     2,
			},
		},
		"test map with arrays": {
			input: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{
						"b": 1,
						"c": "test",
					},
					map[string]interface{}{
						"b": 2,
						"c": "test 2",
					},
					map[string]interface{}{
						"b": 3,
						"c": map[string]interface{}{
							"d": 2,
							"e": "test",
						},
					},
				},
				"f": map[string]interface{}{
					"g": 3,
					"h": "test",
				},
				"i": []interface{}{
					map[string]interface{}{
						"j": 1,
						"k": []interface{}{
							map[string]interface{}{
								"l": 10,
								"m": true,
							},
							map[string]interface{}{
								"l": 20,
								"m": false,
							},
							map[string]interface{}{
								"l": 30,
								"m": true,
							},
						},
					},
					map[string]interface{}{
						"j": 2,
						"k": []interface{}{
							map[string]interface{}{
								"l": 30,
								"m": false,
							},
							map[string]interface{}{
								"l": 20,
								"m": true,
							},
							map[string]interface{}{
								"l": 10,
								"m": false,
							},
						},
					},
				},
				"n": []interface{}{
					1,
					2,
					3,
				},
			},
			expectedOutput: map[string]interface{}{
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
				"n.0":       1,
				"n.1":       2,
				"n.2":       3,
			},
		},
		"test_input_array": {
			input: []interface{}{
				1,
				"abc",
				map[string]interface{}{
					"a": 1,
					"b": 2,
					"c": 3,
				},
				[]interface{}{
					4,
					5,
					"6",
				},
			},
			expectedOutput: map[string]interface{}{
				"0":   1,
				"1":   "abc",
				"2.a": 1,
				"2.b": 2,
				"2.c": 3,
				"3.0": 4,
				"3.1": 5,
				"3.2": "6",
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
