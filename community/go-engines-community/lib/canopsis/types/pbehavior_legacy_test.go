package types_test

import (
	"testing"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	. "github.com/smartystreets/goconvey/convey"
)

func getPBehaviorLegacy() types.PBehaviorLegacy {
	start := types.CpsTime{Time: time.Now()}
	end := types.CpsTime{Time: time.Now().Add(time.Hour * 2)}

	p := types.PBehaviorLegacy{
		Author:        "Alan & Thomas",
		Connector:     "C1995-O1",
		ConnectorName: "space",
		Enabled:       true,
		Filter:        `{"enabled": true}`,
		Name:          "Hale-Bopp",
		Reason:        "Jupiter",
		RRule:         "",
		Start:         &start,
		Stop:          &end,
		Type:          "Maintenance",
	}

	return p
}

func TestPBehaviorLegacy(t *testing.T) {
	Convey("Given a pbehavior and an alarm", t, func() {
		p := getPBehaviorLegacy()
		alarm := getAlarm()

		Convey("Then it should not be impacted with a random alarm", func() {
			So(p.IsImpacted(alarm.EntityID), ShouldBeFalse)
		})

		Convey("Then it should be impacted with a correct alarm", func() {
			p.Eids = []string{alarm.EntityID}
			So(p.IsImpacted(alarm.EntityID), ShouldBeTrue)
		})
	})

	Convey("Given a pbehavior that is active during the first serial of doctor who in 1963", t, func() {
		londonLoc, _ := time.LoadLocation("Europe/London")
		parisLoc, _ := time.LoadLocation("Europe/Paris")
		start := types.CpsTime{Time: time.Date(1963, 11, 23, 17, 15, 0, 0, londonLoc)}
		end := types.CpsTime{Time: time.Date(1963, 11, 23, 17, 40, 0, 0, londonLoc)}
		pbh := types.PBehaviorLegacy{
			Author:        "Verity Lambert",
			Connector:     "The Beeb",
			ConnectorName: "BBC",
			Enabled:       true,
			Filter:        "{}",
			Name:          "Doctor Who - An Unearthly Child",
			Reason:        "The doctor was scared Ian and Barbara will disclose his and his grandaughter presence on earth.",
			RRule:         "FREQ=WEEKLY;BYDAY=SA;COUNT=4",
			Start:         &start,
			Stop:          &end,
			Type:          "Sci-fi",
			Timezone:      "Europe/London",
		}

		Convey("Then it must be active during the broadcast of the serial", func() {
			dateEp1 := time.Date(1963, 11, 23, 17, 30, 0, 0, londonLoc)
			dateEp2 := time.Date(1963, 11, 30, 17, 30, 0, 0, londonLoc)
			dateEp3 := time.Date(1963, 12, 7, 17, 30, 0, 0, londonLoc)
			dateEp4 := time.Date(1963, 12, 14, 17, 30, 0, 0, londonLoc)

			resEp1, errEp1 := pbh.IsActive(dateEp1)
			resEp2, errEp2 := pbh.IsActive(dateEp2)
			resEp3, errEp3 := pbh.IsActive(dateEp3)
			resEp4, errEp4 := pbh.IsActive(dateEp4)

			So(errEp1, ShouldBeNil)
			So(errEp2, ShouldBeNil)
			So(errEp3, ShouldBeNil)
			So(errEp4, ShouldBeNil)

			So(resEp1, ShouldBeTrue)
			So(resEp2, ShouldBeTrue)
			So(resEp3, ShouldBeTrue)
			So(resEp4, ShouldBeTrue)
		})

		Convey("But it must not be active at any time except during the broadcast", func() {
			dateNotAnEp := time.Date(1963, 11, 16, 17, 30, 0, 0, londonLoc)
			resNotAnEp, errNotAnEp := pbh.IsActive(dateNotAnEp)
			So(errNotAnEp, ShouldBeNil)
			So(resNotAnEp, ShouldBeFalse)

			dateNotAnEp = time.Date(2019, 12, 21, 17, 30, 0, 0, londonLoc)
			resNotAnEp, errNotAnEp = pbh.IsActive(dateNotAnEp)
			So(errNotAnEp, ShouldBeNil)
			So(resNotAnEp, ShouldBeFalse)

			dateNotAnEp = time.Date(1963, 12, 23, 17, 30, 0, 0, londonLoc)
			resNotAnEp, errNotAnEp = pbh.IsActive(dateNotAnEp)
			So(errNotAnEp, ShouldBeNil)
			So(resNotAnEp, ShouldBeFalse)
		})

		Convey("Then it must be active for the french viewer at the same time", func() {
			dateEp1 := time.Date(1963, 11, 23, 18, 30, 0, 0, parisLoc)
			dateEp2 := time.Date(1963, 11, 30, 18, 30, 0, 0, parisLoc)
			dateEp3 := time.Date(1963, 12, 7, 18, 30, 0, 0, parisLoc)
			dateEp4 := time.Date(1963, 12, 14, 18, 30, 0, 0, parisLoc)

			resEp1, errEp1 := pbh.IsActive(dateEp1)
			resEp2, errEp2 := pbh.IsActive(dateEp2)
			resEp3, errEp3 := pbh.IsActive(dateEp3)
			resEp4, errEp4 := pbh.IsActive(dateEp4)

			So(errEp1, ShouldBeNil)
			So(errEp2, ShouldBeNil)
			So(errEp3, ShouldBeNil)
			So(errEp4, ShouldBeNil)

			So(resEp1, ShouldBeTrue)
			So(resEp2, ShouldBeTrue)
			So(resEp3, ShouldBeTrue)
			So(resEp4, ShouldBeTrue)
		})

		Convey("But it must not be active for the french viewer at any time except during the broadcast", func() {
			dateNotAnEp := time.Date(1963, 12, 21, 16, 30, 0, 0, parisLoc)
			resNotAnEp, errNotAnEp := pbh.IsActive(dateNotAnEp)
			So(errNotAnEp, ShouldBeNil)
			So(resNotAnEp, ShouldBeFalse)

			dateNotAnEp = time.Date(2019, 12, 21, 23, 30, 0, 0, parisLoc)
			resNotAnEp, errNotAnEp = pbh.IsActive(dateNotAnEp)
			So(errNotAnEp, ShouldBeNil)
			So(resNotAnEp, ShouldBeFalse)

			dateNotAnEp = time.Date(1963, 11, 16, 16, 30, 0, 0, parisLoc)
			resNotAnEp, errNotAnEp = pbh.IsActive(dateNotAnEp)
			So(errNotAnEp, ShouldBeNil)
			So(resNotAnEp, ShouldBeFalse)
		})
	})

	Convey("Given a pbehavior that represent the missing episode of doctor who serial `The Reign of Terror`", t, func() {
		londonLoc, _ := time.LoadLocation("Europe/London")
		parisLoc, _ := time.LoadLocation("Europe/Paris")
		start := types.CpsTime{Time: time.Date(1967, 1, 14, 17, 15, 0, 0, londonLoc)}
		end := types.CpsTime{Time: time.Date(1967, 1, 14, 17, 40, 0, 0, londonLoc)}
		exdateEp4 := types.CpsTime{time.Date(1967, 1, 14, 17, 15, 0, 0, londonLoc)}
		exdateEp5 := types.CpsTime{time.Date(1967, 2, 4, 17, 15, 0, 0, londonLoc)}
		pbh := types.PBehaviorLegacy{
			Author:        "Verity Lambert",
			Connector:     "The Beeb",
			ConnectorName: "BBC",
			Enabled:       true,
			Filter:        "{}",
			Name:          "Doctor Who - The Daleks' Master Plan",
			Reason:        "in the year 4000 the daleks try to conquer the solar system with a weapon that can destroy the fabric of time.",
			RRule:         "FREQ=WEEKLY;BYDAY=SA;COUNT=60",
			Start:         &start,
			Stop:          &end,
			Type:          "Sci-fi",
			Timezone:      "Europe/London",
			Exdate:        []types.CpsTime{exdateEp4, exdateEp5},
		}

		Convey("Then it should only be active for during the non missing episode", func() {
			dateEp1 := time.Date(1967, 1, 14, 17, 30, 0, 0, londonLoc)
			dateEp2 := time.Date(1967, 1, 21, 17, 30, 0, 0, londonLoc)
			dateEp3 := time.Date(1967, 1, 28, 17, 30, 0, 0, londonLoc)
			dateEp4 := time.Date(1967, 2, 3, 17, 30, 0, 0, londonLoc)

			resEp1, errEp1 := pbh.IsActive(dateEp1)
			resEp2, errEp2 := pbh.IsActive(dateEp2)
			resEp3, errEp3 := pbh.IsActive(dateEp3)
			resEp4, errEp4 := pbh.IsActive(dateEp4)

			So(errEp1, ShouldBeNil)
			So(errEp2, ShouldBeNil)
			So(errEp3, ShouldBeNil)
			So(errEp4, ShouldBeNil)

			So(resEp1, ShouldBeFalse)
			So(resEp2, ShouldBeTrue)
			So(resEp3, ShouldBeTrue)
			So(resEp4, ShouldBeFalse)
		})

		Convey("Then we must have the same result for french viewer.", func() {
			dateEp1 := time.Date(1967, 1, 14, 18, 30, 0, 0, parisLoc)
			dateEp2 := time.Date(1967, 1, 21, 18, 30, 0, 0, parisLoc)
			dateEp3 := time.Date(1967, 1, 28, 18, 30, 0, 0, parisLoc)
			dateEp4 := time.Date(1967, 2, 4, 18, 30, 0, 0, parisLoc)

			resEp1, errEp1 := pbh.IsActive(dateEp1)
			resEp2, errEp2 := pbh.IsActive(dateEp2)
			resEp3, errEp3 := pbh.IsActive(dateEp3)
			resEp4, errEp4 := pbh.IsActive(dateEp4)

			So(errEp1, ShouldBeNil)
			So(errEp2, ShouldBeNil)
			So(errEp3, ShouldBeNil)
			So(errEp4, ShouldBeNil)

			So(resEp1, ShouldBeFalse)
			So(resEp2, ShouldBeTrue)
			So(resEp3, ShouldBeTrue)
			So(resEp4, ShouldBeFalse)
		})

		Convey("Then we must have the same result if the missing date match the french broadcast hour.", func() {
			frExdateEp4 := types.CpsTime{time.Date(1967, 1, 14, 18, 15, 0, 0, parisLoc)}
			frExdateEp5 := types.CpsTime{time.Date(1967, 2, 4, 18, 15, 0, 0, parisLoc)}
			pbh.Exdate = []types.CpsTime{frExdateEp4, frExdateEp5}

			dateEp1 := time.Date(1967, 1, 14, 18, 30, 0, 0, parisLoc)
			dateEp2 := time.Date(1967, 1, 21, 18, 30, 0, 0, parisLoc)
			dateEp3 := time.Date(1967, 1, 28, 18, 30, 0, 0, parisLoc)
			dateEp4 := time.Date(1967, 2, 4, 18, 30, 0, 0, parisLoc)

			resEp1, errEp1 := pbh.IsActive(dateEp1)
			resEp2, errEp2 := pbh.IsActive(dateEp2)
			resEp3, errEp3 := pbh.IsActive(dateEp3)
			resEp4, errEp4 := pbh.IsActive(dateEp4)

			So(errEp1, ShouldBeNil)
			So(errEp2, ShouldBeNil)
			So(errEp3, ShouldBeNil)
			So(errEp4, ShouldBeNil)

			So(resEp1, ShouldBeFalse)
			So(resEp2, ShouldBeTrue)
			So(resEp3, ShouldBeTrue)
			So(resEp4, ShouldBeFalse)
		})
	})
}
