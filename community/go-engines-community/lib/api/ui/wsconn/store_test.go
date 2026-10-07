package wsconn

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestFreshnessFilterUsesNumericTimestamp(t *testing.T) {
	store := &store{readInterval: time.Minute}
	cutoff := store.freshnessFilter()["updated"].(bson.M)["$gt"]
	cutoffType, _, err := bson.MarshalValue(cutoff)
	if err != nil {
		t.Fatalf("marshal freshness cutoff: %v", err)
	}
	if cutoffType != bson.TypeInt64 {
		t.Fatalf("freshness cutoff BSON type = %s, want %s", cutoffType, bson.TypeInt64)
	}
}
