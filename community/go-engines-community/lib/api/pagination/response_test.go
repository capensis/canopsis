package pagination

import (
	"fmt"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestNewMeta(t *testing.T) {
	testCases := []struct {
		query    Query
		total    int64
		expected MetaResponse
	}{
		{
			query: Query{
				Page:     1,
				Limit:    10,
				Paginate: true,
			},
			total: 34,
			expected: MetaResponse{
				Page:       1,
				PerPage:    10,
				PageCount:  4,
				TotalCount: 34,
			},
		},
		{
			query: Query{
				Page:     1,
				Limit:    10,
				Paginate: true,
			},
			total: 3,
			expected: MetaResponse{
				Page:       1,
				PerPage:    10,
				PageCount:  1,
				TotalCount: 3,
			},
		},
		{
			query: Query{
				Page:     1,
				Limit:    0,
				Paginate: true,
			},
			total: 0,
			expected: MetaResponse{
				Page:       1,
				PerPage:    0,
				PageCount:  1,
				TotalCount: 0,
			},
		},
		{
			query: Query{
				Page:     1,
				Limit:    10,
				Paginate: true,
			},
			total: 0,
			expected: MetaResponse{
				Page:       1,
				PerPage:    10,
				PageCount:  1,
				TotalCount: 0,
			},
		},
		{
			query: Query{
				Page:     1,
				Limit:    0,
				Paginate: true,
			},
			total: 34,
			expected: MetaResponse{
				Page:       1,
				PerPage:    0,
				PageCount:  1,
				TotalCount: 34,
			},
		},
		{
			query: Query{
				Page:     1,
				Limit:    10,
				Paginate: false,
			},
			total: 34,
			expected: MetaResponse{
				Page:       1,
				PerPage:    34,
				PageCount:  1,
				TotalCount: 34,
			},
		},
		{
			query: Query{
				Page:     1,
				Limit:    0,
				Paginate: false,
			},
			total: 34,
			expected: MetaResponse{
				Page:       1,
				PerPage:    34,
				PageCount:  1,
				TotalCount: 34,
			},
		},
		{
			query: Query{
				Page:     1,
				Limit:    10,
				Paginate: false,
			},
			total: 0,
			expected: MetaResponse{
				Page:       1,
				PerPage:    0,
				PageCount:  1,
				TotalCount: 0,
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			received := NewMeta(tc.query, tc.total)
			if diff := pretty.Compare(tc.expected, received); diff != "" {
				t.Errorf("diff: (-want +got)\n%s", diff)
			}
		})
	}
}
