package alarm

import (
	"regexp"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/fielddoc"
)

var (
	dynamicFieldSamples = []string{"infos.rule.name", "entity.infos.site.value", "entity.component_infos.site.value"}
	invalidFieldSamples = []string{"", "foo", "infos.rule", "infos.a.b.c", "entity.infos.site.name", "entity.foo.site.value", "v.display_name"}
)

func TestSortableFields_GivenDocumentedField_ShouldBeAcceptedByResolveSortableField(t *testing.T) {
	assertFieldsDocMatchesResolver(t, SortableFields(), func(f, prefix string) bool {
		_, _, ok := ResolveSortableField(f, prefix)

		return ok
	})
}

func TestSearchableFields_GivenDocumentedField_ShouldBeAcceptedByResolveSearchField(t *testing.T) {
	assertFieldsDocMatchesResolver(t, SearchableFields(), func(f, prefix string) bool {
		_, _, ok := resolveSearchField(f, prefix)

		return ok
	})
}

func TestSortableFields_SortRegexp(t *testing.T) {
	re := regexp.MustCompile(SortableFields().SortRegexp())
	for _, s := range []string{"creation_date", "creation_date:desc", "creation_date:desc,entity.name,infos.a.b:asc", "entity.component_infos.x.value:desc"} {
		if !re.MatchString(s) {
			t.Errorf("expected %q to match", s)
		}
	}

	for _, s := range []string{"", "foo:desc", "creation_date:up", "creation_date,", ",creation_date", "creation_date;name"} {
		if re.MatchString(s) {
			t.Errorf("expected %q not to match", s)
		}
	}
}

func assertFieldsDocMatchesResolver(t *testing.T, doc fielddoc.Fields, accepts func(f, prefix string) bool) {
	t.Helper()
	field := regexp.MustCompile("^" + doc.Regexp() + "$")
	for _, prefix := range []string{EntityPrefix(ResolutionOpen), EntityPrefix(ResolutionResolved)} {
		for _, f := range append(doc.Static, dynamicFieldSamples...) {
			if !field.MatchString(f) {
				t.Errorf("expected %q to be documented", f)
			}

			if !accepts(f, prefix) {
				t.Errorf("documented field %q is not accepted with prefix %q", f, prefix)
			}
		}

		for _, f := range invalidFieldSamples {
			if field.MatchString(f) {
				t.Errorf("expected %q not to be documented", f)
			}

			if accepts(f, prefix) {
				t.Errorf("undocumented field %q is accepted with prefix %q", f, prefix)
			}
		}
	}
}
