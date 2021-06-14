package testutils

import (
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/config"
	"os"
	"testing"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Testutils constants
const (
	EnvCpsTestSkipLong = "CPS_TEST_SKIP_LONG"
)

// EnvBackup handles backup and restore of environment variables.
// 	{
// 		e := NewEnvBackup("MY_VAR", "newval")
//		doYourStuff()
// 		e.Restore()
// 	}
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

// MongoDrop handles collection drop even if the collection doesn't exists.
func MongoDrop(collection *mgo.Collection) error {
	if err := collection.Insert(bson.M{"_id": "blank"}); err != nil {
		return fmt.Errorf("mongo drop: insert: %v", err)
	}

	if err := collection.DropCollection(); err != nil {
		return fmt.Errorf("mongo drop: drop: %v", err)
	}

	if err := collection.Create(&mgo.CollectionInfo{}); err != nil {
		return fmt.Errorf("mongo create: %v", err)
	}

	return nil
}

func GetTestConf() config.CanopsisConf {
	return config.CanopsisConf{
		Global: config.SectionGlobal{
			PrefetchCount: 10000,
			PrefetchSize:  0,
		},
		Alarm: config.SectionAlarm{
			FlappingFreqLimit:    1,
			FlappingInterval:     1,
			StealthyInterval:     100,
			CancelAutosolveDelay: "1h",
		},
	}
}