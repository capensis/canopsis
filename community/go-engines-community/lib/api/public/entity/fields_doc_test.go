package entity

import (
	"regexp"
	"testing"
)

func TestFieldsDoc_GivenField_ShouldMatchValidation(t *testing.T) {
	valid := []string{"infos.site.value", "component_infos.site.value"}
	invalid := []string{"", "foo", "infos.site", "infos.site.name", "entity.infos.site.value"}
	for name, tc := range map[string]struct {
		re      *regexp.Regexp
		static  []string
		accepts func(string) bool
	}{
		"sort":   {regexp.MustCompile("^" + SortableFields().Regexp() + "$"), SortableFields().Static, IsSortableField},
		"search": {regexp.MustCompile("^" + SearchableFields().Regexp() + "$"), SearchableFields().Static, isSearchField},
	} {
		t.Run(name, func(t *testing.T) {
			for _, f := range append(tc.static, valid...) {
				if !tc.re.MatchString(f) || !tc.accepts(f) {
					t.Errorf("%q: documented=%v accepted=%v, expected both", f, tc.re.MatchString(f), tc.accepts(f))
				}
			}

			for _, f := range invalid {
				if tc.re.MatchString(f) || tc.accepts(f) {
					t.Errorf("%q: documented=%v accepted=%v, expected neither", f, tc.re.MatchString(f), tc.accepts(f))
				}
			}
		})
	}
}
