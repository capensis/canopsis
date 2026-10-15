package pagination

import (
	"fmt"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestNewResponse(t *testing.T) {
	type Item struct{ ID string }
	testCases := []struct {
		query    Query
		data     []Item
		expected ListResponse[Item]
	}{
		{
			query: Query{
				Offset: 0,
				Limit:  5,
			},
			data: []Item{{ID: "1"}, {ID: "2"}, {ID: "3"}, {ID: "4"}, {ID: "5"}},
			expected: ListResponse[Item]{
				Data: []Item{{ID: "1"}, {ID: "2"}, {ID: "3"}, {ID: "4"}, {ID: "5"}},
				Meta: MetaResponse{
					Offset:  0,
					Limit:   5,
					HasMore: false,
				},
			},
		},
		{
			query: Query{
				Offset: 0,
				Limit:  5,
			},
			data: []Item{{ID: "1"}, {ID: "2"}, {ID: "3"}, {ID: "4"}, {ID: "5"}, {ID: "6"}},
			expected: ListResponse[Item]{
				Data: []Item{{ID: "1"}, {ID: "2"}, {ID: "3"}, {ID: "4"}, {ID: "5"}},
				Meta: MetaResponse{
					Offset:  0,
					Limit:   5,
					HasMore: true,
				},
			},
		},
		{
			query: Query{
				Offset: 5,
				Limit:  5,
			},
			data: []Item{{ID: "6"}, {ID: "7"}, {ID: "8"}, {ID: "9"}, {ID: "10"}, {ID: "11"}},
			expected: ListResponse[Item]{
				Data: []Item{{ID: "6"}, {ID: "7"}, {ID: "8"}, {ID: "9"}, {ID: "10"}},
				Meta: MetaResponse{
					Offset:  5,
					Limit:   5,
					HasMore: true,
				},
			},
		},
		{
			query: Query{
				Offset: 0,
				Limit:  5,
			},
			data: []Item{},
			expected: ListResponse[Item]{
				Data: []Item{},
				Meta: MetaResponse{
					Offset:  0,
					Limit:   5,
					HasMore: false,
				},
			},
		},
		{
			query: Query{
				Offset: 10,
				Limit:  5,
			},
			data: []Item{},
			expected: ListResponse[Item]{
				Data: []Item{},
				Meta: MetaResponse{
					Offset:  10,
					Limit:   5,
					HasMore: false,
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			received := NewResponse(tc.query, tc.data)
			if diff := pretty.Compare(tc.expected, received); diff != "" {
				t.Errorf("diff: (-want +got)\n%s", diff)
			}
		})
	}
}
