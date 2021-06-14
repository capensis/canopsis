package bulk

import (
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// BulkSizeMax is the maximum number of bulk operations that can be handled before a flush
// to database is triggered.
const BulkSizeMax = 1000

// BulkPacketSize maximum document size for MongoDB is 16793600 bytes. To avoid reaching
// that value, set max size to a bit less to avoid any overhead.
const BulkPacketSize = 16000000

type unsafe struct {
	collection  *mgo.Collection
	ninserts    int // track currently inserted operations to handle flush if required
	nupdates    int
	inserts     *mgo.Bulk
	updates     *mgo.Bulk
	uniqupdates map[string]OpUpdate
	usize       uint64
	isize       uint64
	max         int
}

// New creates an thread-UNSAFE bulk struct and returns it under the Bulk interface.
// Insert operations are unordered while updates are ordered.
//
// max is used to limit the number of documents you can queue in a bulk before it is
// flushed to the database.
//
// If max > BulkSizeMax it's value is limited to BulkSizeMax.
func New(collection *mgo.Collection, max int) Bulk {
	if max > BulkSizeMax {
		max = BulkSizeMax
	}

	ibulk := collection.Bulk()
	ubulk := collection.Bulk()

	ibulk.Unordered()

	ab := unsafe{
		collection:  collection,
		inserts:     ibulk,
		updates:     ubulk,
		max:         max,
		uniqupdates: make(map[string]OpUpdate),
	}

	return &ab
}

func (a *unsafe) AddInsert(insert OpInsert) (*mgo.BulkResult, error) {
	csize := uint64(0)
	b, _ := bson.Marshal(insert.Data())
	csize += uint64(len(b))

	a.isize += csize

	if a.ninserts == a.max || a.isize > BulkPacketSize {
		br, err := a.PerformInserts()
		if err != nil {
			return br, fmt.Errorf("bulk add insert: %v", err)
		}
	}

	a.inserts.Insert(insert.Data())
	a.ninserts++

	return nil, nil
}

func (a *unsafe) AddUpdate(update OpUpdate) (*mgo.BulkResult, error) {
	csize := uint64(0)
	b, _ := bson.Marshal(update.Selector())
	csize += uint64(len(b))
	b, _ = bson.Marshal(update.Data())
	csize += uint64(len(b))

	a.usize += csize

	if a.nupdates == a.max || a.usize > BulkPacketSize {
		br, err := a.PerformUpdates()
		if err != nil {
			return br, fmt.Errorf("bulk add update: %v", err)
		}
	}

	a.updates.Update(update.Selector(), update.Data())
	a.nupdates++

	return nil, nil
}

func (a *unsafe) AddUpsert(update OpUpdate) (*mgo.BulkResult, error) {
	csize := uint64(0)
	b, _ := bson.Marshal(update.Selector())
	csize += uint64(len(b))
	b, _ = bson.Marshal(update.Data())
	csize += uint64(len(b))

	a.usize += csize

	if a.nupdates == a.max || a.usize > BulkPacketSize {
		br, err := a.PerformUpdates()
		if err != nil {
			return br, fmt.Errorf("bulk add update: %v", err)
		}
	}

	a.updates.Upsert(update.Selector(), update.Data())
	a.nupdates++

	return nil, nil
}

func (a *unsafe) AddUniqUpdate(id string, update OpUpdate) (*mgo.BulkResult, error) {
	a.uniqupdates[id] = update

	if len(a.uniqupdates) > a.max {
		br, err := a.PerformUniqUpdates()
		if err != nil {
			return br, fmt.Errorf("bulk add uniq update: %v", err)
		}
		return br, nil
	}

	return nil, nil
}

func (a *unsafe) PerformInserts() (*mgo.BulkResult, error) {
	a.collection.Database.Session.Refresh()

	br, err := a.inserts.Run()
	a.ninserts = 0
	a.inserts = a.collection.Bulk()
	a.isize = 0
	a.inserts.Unordered()
	return br, err
}

func (a *unsafe) PerformUpdates() (*mgo.BulkResult, error) {
	a.collection.Database.Session.Refresh()

	br, err := a.updates.Run()
	a.nupdates = 0
	a.usize = 0
	a.updates = a.collection.Bulk()
	return br, err
}

func (a *unsafe) resetUniqUpdates() {
	a.uniqupdates = make(map[string]OpUpdate)
}

func (a *unsafe) PerformUniqUpdates() (*mgo.BulkResult, error) {
	nmodified := 0
	nmatched := 0

	for _, op := range a.uniqupdates {
		br, err := a.AddUpdate(op)

		if br != nil {
			nmatched += br.Matched
			nmodified += br.Modified
		}

		if err != nil {
			br.Matched += nmatched
			br.Modified += nmodified
			a.resetUniqUpdates()
			return br, fmt.Errorf("bulk perform uniq updates: %v", err)
		}
	}

	br, err := a.PerformUpdates()
	if err != nil {
		a.resetUniqUpdates()
		return br, fmt.Errorf("bulk perform uniq updates: %v", err)
	}

	if br != nil {
		nmatched += br.Matched
		nmodified += br.Modified
	}

	bulkRes := mgo.BulkResult{
		Modified: nmodified,
		Matched:  nmatched,
	}

	a.resetUniqUpdates()
	return &bulkRes, nil
}

func (a *unsafe) Close() {
	a.PerformInserts()
	a.PerformUpdates()
}
