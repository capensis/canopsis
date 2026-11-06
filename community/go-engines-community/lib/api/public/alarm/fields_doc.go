package alarm

import (
	"slices"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/fielddoc"
)

// Patterns must stay in sync with isInfoValueField and resolveEntityInfoValueField.
var (
	infoValueFieldPattern       = `infos\.[^.]+\.[^.]+`
	entityInfoValueFieldPattern = `entity\.(?:infos|component_infos)\.[^.]+\.value`
)

// SortableFields describes the fields accepted by ResolveSortableField.
func SortableFields() fielddoc.Fields {
	return fielddoc.Fields{
		Static:   slices.Concat(sortableFields, entitySortableFields),
		Labels:   strings.Fields(SortErrParam()),
		Patterns: []string{infoValueFieldPattern, entityInfoValueFieldPattern},
	}
}

// SearchableFields describes the fields accepted by resolveSearchField.
func SearchableFields() fielddoc.Fields {
	return fielddoc.Fields{
		Static:   slices.Concat(searchableFields, entitySearchableFields),
		Labels:   strings.Fields(searchErrParam()),
		Patterns: []string{infoValueFieldPattern, entityInfoValueFieldPattern},
	}
}

// StepSortableFields describes the fields accepted by the steps sort.
func StepSortableFields() fielddoc.Fields {
	return fielddoc.Fields{
		Static: slices.Clone(stepSortableFields),
		Labels: slices.Clone(stepSortableFields),
	}
}
