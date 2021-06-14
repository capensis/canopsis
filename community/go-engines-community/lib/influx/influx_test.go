package influx_test

import (
	"os"
	"testing"

	"git.canopsis.net/canopsis/go-engines/lib/influx"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewInfluxSession(t *testing.T) {
	Convey("Wanna check some good InfluxDB session?", t, func() {
		Convey("Bad url", func() {
			ou := os.Getenv(influx.EnvURL)
			os.Setenv(influx.EnvURL, "-... .- -.. / ..- .-. .-..")
			_, err := influx.NewSession()
			So(err, ShouldNotBeNil)
			os.Setenv(influx.EnvURL, ou)
		})

		Convey("all good", func() {
			_, err := influx.NewSession()
			So(err, ShouldBeNil)
		})
	})
}
