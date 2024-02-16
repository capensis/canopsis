package testutils

import (
	"os"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
)

// Testutils constants
const (
	EnvCpsTestSkipLong = "CPS_TEST_SKIP_LONG"
)

// EnvBackup handles backup and restore of environment variables.
//
//	{
//		e := NewEnvBackup("MY_VAR", "newval")
//		doYourStuff()
//		e.Restore()
//	}
type EnvBackup struct {
	value string
	name  string
}

// NewEnvBackup create EnvBackup, then run SaveSet(newval) before returning.
func NewEnvBackup(name string, newval string) (EnvBackup, error) {
	e := EnvBackup{
		name: name,
	}
	err := e.SaveSet(newval)
	return e, err
}

// SaveSet backup the current value of EnvBackup.name
// then set the env var to newval.
func (e *EnvBackup) SaveSet(newval string) error {
	e.value = os.Getenv(e.name)
	return os.Setenv(e.name, newval)
}

// Restore set the env var to it's backuped value.
func (e EnvBackup) Restore() error {
	return os.Setenv(e.name, e.value)
}

// SkipLongIfSet set the current test as being a "long" running one.
// If CPS_TEST_SKIP_LONG exactly equals to "1" this test is testing.T.SkipNow()
func SkipLongIfSet(t *testing.T) {
	if os.Getenv(EnvCpsTestSkipLong) == "1" {
		t.Skipf("skipped long test %s", t.Name())
	}
}

func GetTestConf() config.CanopsisConf {
	return config.CanopsisConf{
		Global: config.SectionGlobal{
			PrefetchCount: 10000,
			PrefetchSize:  0,
		},
		Alarm: config.SectionAlarm{
			StealthyInterval:     100,
			CancelAutosolveDelay: "1h",
		},
	}
}
