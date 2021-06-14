package bulk

import (
	"github.com/globalsign/mgo"
)

// Bulk allows you to send documents to MongoDB by batches.
// See bulk.New() and bulk.NewSafe() for more details.
type Bulk interface {
	// AddInsert handle up to max documents per add
	AddInsert(insert OpInsert) (*mgo.BulkResult, error)
	// AddUpdate handle up to max documents per add
	AddUpdate(update OpUpdate) (*mgo.BulkResult, error)

	AddUpsert(update OpUpdate) (*mgo.BulkResult, error)

	// AddUniqUpdate overrides updates on the same id as long as the bulk
	// isnt flushed.
	AddUniqUpdate(id string, update OpUpdate) (*mgo.BulkResult, error)

	// PerformInserts run the bulk operation
	PerformInserts() (*mgo.BulkResult, error)
	// PerformUpdates run the bulk operation
	PerformUpdates() (*mgo.BulkResult, error)

	// PerformUniqUpdates ...
	PerformUniqUpdates() (*mgo.BulkResult, error)
}
