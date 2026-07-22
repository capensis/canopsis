package mongoquery

import (
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestGetSearchQuery_GivenValidRegex_ShouldKeepPattern(t *testing.T) {
	query := GetSearchQuery("foo.*bar", []string{"name"})

	actual := query["$or"].([]bson.M)
	regex := actual[0]["name"].(bson.Regex)
	if regex.Pattern != "foo.*bar" {
		t.Fatalf("expected regex pattern %q but got %q", "foo.*bar", regex.Pattern)
	}
	if regex.Options != "i" {
		t.Fatalf("expected regex option %q but got %q", "i", regex.Options)
	}
}

func TestGetSearchQuery_GivenInvalidRegex_ShouldEscapePattern(t *testing.T) {
	query := GetSearchQuery("[", []string{"name"})

	actual := query["$or"].([]bson.M)
	regex := actual[0]["name"].(bson.Regex)
	if regex.Pattern != ".*\\[.*" {
		t.Fatalf("expected escaped regex pattern %q but got %q", ".*\\[.*", regex.Pattern)
	}
	if regex.Options != "i" {
		t.Fatalf("expected regex option %q but got %q", "i", regex.Options)
	}
}
