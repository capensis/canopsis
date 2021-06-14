package bulk

import (
	"sync"

	"github.com/globalsign/mgo"
)

type safe struct {
	b Bulk
	l sync.RWMutex
}

// NewSafe wrap the given unsafe Bulk to a thread safe bulk using sync.RWMutex.
func NewSafe(unsafe Bulk) Bulk {
	b := safe{
		b: unsafe,
	}
	return &b
}

func (b *safe) AddInsert(insert OpInsert) (*mgo.BulkResult, error) {
	b.l.Lock()
	defer b.l.Unlock()
	return b.b.AddInsert(insert)
}

func (b *safe) AddUpdate(update OpUpdate) (*mgo.BulkResult, error) {
	b.l.Lock()
	defer b.l.Unlock()
	return b.b.AddUpdate(update)
}

func (b *safe) AddUpsert(update OpUpdate) (*mgo.BulkResult, error) {
	b.l.Lock()
	defer b.l.Unlock()
	return b.b.AddUpsert(update)
}

func (b *safe) AddUniqUpdate(id string, update OpUpdate) (*mgo.BulkResult, error) {
	b.l.Lock()
	defer b.l.Unlock()
	return b.b.AddUniqUpdate(id, update)
}

func (b *safe) PerformInserts() (*mgo.BulkResult, error) {
	b.l.Lock()
	defer b.l.Unlock()
	return b.b.PerformInserts()
}

func (b *safe) PerformUpdates() (*mgo.BulkResult, error) {
	b.l.Lock()
	defer b.l.Unlock()
	return b.b.PerformUpdates()
}

func (b *safe) PerformUniqUpdates() (*mgo.BulkResult, error) {
	b.l.Lock()
	defer b.l.Unlock()
	return b.b.PerformUniqUpdates()
}
