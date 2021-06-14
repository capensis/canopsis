package context

import "strings"

// EnrichFields is a simple include/exclude ruleset.
type EnrichFields struct {
	includes     map[string]bool
	someIncludes bool
	excludes     map[string]bool
	someExcludes bool
}

// Allow returns true if the given field
// is part of included fields or not in excluded fields.
func (ef *EnrichFields) Allow(field string) bool {
	if ef.someIncludes {
		if _, ok := ef.includes[field]; ok {
			return true
		}
		return false
	}

	if ef.someExcludes {
		if _, ok := ef.excludes[field]; ok {
			return false
		}
		return true
	}

	return true
}

// AddInclude ...
func (ef *EnrichFields) AddInclude(field string) {
	ef.includes[field] = true
	ef.someIncludes = true
}

// AddExclude ...
func (ef *EnrichFields) AddExclude(field string) {
	ef.excludes[field] = true
	ef.someExcludes = true
}

// DelInclude ...
func (ef *EnrichFields) DelInclude(field string) {
	delete(ef.includes, field)
	if len(ef.includes) == 0 {
		ef.someIncludes = false
	}
}

// DelExclude ...
func (ef *EnrichFields) DelExclude(field string) {
	delete(ef.excludes, field)
	if len(ef.includes) == 0 {
		ef.someExcludes = false
	}
}

// NewEnrichFields only works on include if not empty, otherwise exclude.
// include and exclude must be strings of the form:
//
//    "field1,field2,field3"
//
// Any field containing a coma is not supported through the constructor,
// use AddInclude/AddExclude manually.
func NewEnrichFields(include, exclude string) EnrichFields {
	ef := EnrichFields{
		includes:     make(map[string]bool),
		excludes:     make(map[string]bool),
		someIncludes: false,
		someExcludes: false,
	}
	excludes := strings.Split(exclude, ",")
	includes := strings.Split(include, ",")

	if len(includes) > 0 && includes[0] != "" {
		for _, field := range includes {
			ef.AddInclude(field)
		}

		return ef
	}

	if len(excludes) > 0 && excludes[0] != "" {
		for _, field := range excludes {
			ef.AddExclude(field)
		}
	}

	return ef
}
