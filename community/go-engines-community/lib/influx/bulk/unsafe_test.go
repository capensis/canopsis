package bulk_test

import (
	"testing"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/influx"
	influxmod "github.com/influxdata/influxdb/client/v2"

	"git.canopsis.net/canopsis/go-engines/lib/influx/bulk"
	. "github.com/smartystreets/goconvey/convey"
)

func TestInfluxBulk(t *testing.T) {
	Convey("Setup", t, func() {
		client, err := influx.NewSession()
		So(err, ShouldBeNil)

		config := influxmod.BatchPointsConfig{
			Database: "canopsis",
		}

		b, err := bulk.New(2, client, config)
		So(err, ShouldBeNil)
		So(b, ShouldNotBeNil)

		Convey("Setup point", func() {
			fields := make(map[string]interface{})
			fields["f1"] = 1

			tags := make(map[string]string)
			tags["t1"] = "t1"

			po := bulk.PointOp{
				Fields: fields,
				Tags:   tags,
				Name:   "p1",
				Time:   time.Now(),
			}

			r, err := client.Query(influxmod.NewQuery("DELETE FROM p1", "canopsis", "s"))
			So(err, ShouldBeNil)
			So(r.Error(), ShouldBeNil)

			Convey("Add the first (1) point", func() {

				So(b.AddPoints(po), ShouldBeNil)
				po.Time = po.Time.Add(time.Second)

				r, err = client.Query(influxmod.NewQuery("SELECT * FROM p1", "canopsis", "s"))
				So(err, ShouldBeNil)
				So(len(r.Results), ShouldEqual, 1)
				So(len(r.Results[0].Series), ShouldEqual, 0)

				Convey("Add second (2) point", func() {
					po.Time = po.Time.Add(time.Second)
					So(b.AddPoints(po), ShouldBeNil)

					r, err = client.Query(influxmod.NewQuery("SELECT * FROM p1", "canopsis", "s"))
					So(err, ShouldBeNil)
					So(len(r.Results), ShouldEqual, 1)
					So(len(r.Results[0].Series), ShouldEqual, 0)

					Convey("Add third (3) point to auto perform", func() {
						po.Time = po.Time.Add(time.Second)
						So(b.AddPoints(po), ShouldBeNil)

						r, err = client.Query(influxmod.NewQuery("SELECT * FROM p1", "canopsis", "s"))
						So(err, ShouldBeNil)
						So(len(r.Results), ShouldEqual, 1)
						So(len(r.Results[0].Series), ShouldEqual, 1)
						So(len(r.Results[0].Series[0].Values), ShouldEqual, 3)

						Convey("Add one point and force perform", func() {
							po.Time = po.Time.Add(time.Second)
							So(b.AddPoints(po), ShouldBeNil)
							So(b.Perform(), ShouldBeNil)

							r, err = client.Query(influxmod.NewQuery("SELECT * FROM p1", "canopsis", "s"))
							So(err, ShouldBeNil)
							So(len(r.Results), ShouldEqual, 1)
							So(len(r.Results[0].Series), ShouldEqual, 1)
							So(len(r.Results[0].Series[0].Values), ShouldEqual, 4)
						})
					})
				})
			})
		})
	})
}
