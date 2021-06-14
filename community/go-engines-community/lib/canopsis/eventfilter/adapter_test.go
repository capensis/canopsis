package eventfilter_test

import (
	"testing"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"github.com/globalsign/mgo/bson"
	. "github.com/smartystreets/goconvey/convey"
)

func testNewAdapter(rules ...bson.M) eventfilter.Adapter {
	session, err := mongo.NewSession(mongo.Timeout)
	So(err, ShouldBeNil)

	collection := eventfilter.DefaultCollection(session)

	_, err = collection.RemoveAll()
	So(err, ShouldBeNil)

	// Add rules to the collection
	for _, rule := range rules {
		So(collection.Insert(rule), ShouldBeNil)
	}

	return eventfilter.NewAdapter(collection)
}

func TestList(t *testing.T) {
	Convey("Given a filled rules collection", t, func() {
		adapter := testNewAdapter(
			untypedRule,
			invalidTypeRule,
			mistypedTypeRule,
			mistypedPatternRule,
			invalidRegexPatternRule,
			unexpectedFieldPatternRule,
			unexpectedFieldRule,
			unexpectedEnrichmentFieldRule,
			noActionsRule,
			emptyActionsRule,
			invalidOutcomeRule,
			mistypedOutcomeRule,
			invalidActionRule,
			missingActionTypeRule,
			unexpectedActionFieldRule,
			missingActionFieldRule,
			disabledRule,
			dropRule,
			breakRule,
			enrichmentRule,
			failingEnrichmentRule,
			translationRule,
		)

		Convey("Listing the rules returns all the rules", func() {
			rules, err := adapter.List()

			So(err, ShouldBeNil)

			Convey("Only the valid rules should be returned", func() {
				So(rules, ShouldHaveLength, 5)

				So(rules[0].ID, ShouldEqual, "valid_break")
				So(rules[1].ID, ShouldEqual, "valid_enrichment")
				So(rules[2].ID, ShouldEqual, "valid_drop")
				So(rules[3].ID, ShouldEqual, "failing_enrichment")
				So(rules[4].ID, ShouldEqual, "translation")
			})
		})
	})
}
