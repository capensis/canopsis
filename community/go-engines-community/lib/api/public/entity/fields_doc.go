package entity

import (
	"slices"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/fielddoc"
)

// Must stay in sync with isInfoValueField.
var infoValueFieldPattern = `(?:infos|component_infos)\.[^.]+\.value`

// SortableFields describes the fields accepted by IsSortableField.
func SortableFields() fielddoc.Fields {
	return fielddoc.Fields{
		Static:   slices.Clone(sortableFields),
		Labels:   strings.Fields(SortErrParam()),
		Patterns: []string{infoValueFieldPattern},
	}
}

// SearchableFields describes the fields accepted by isSearchField.
func SearchableFields() fielddoc.Fields {
	return fielddoc.Fields{
		Static:   slices.Clone(searchableFields),
		Labels:   strings.Fields(searchErrParam()),
		Patterns: []string{infoValueFieldPattern},
	}
}
