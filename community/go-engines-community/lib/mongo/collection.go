package mongo

import (
	"log"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/mongo/bulk"
	"github.com/globalsign/mgo"
)

const (
	// maxRetries is the maximum number of attempts made to resend a query to
	// the MongoDB database when the connection fails.
	maxRetries int = 5
)

// retryDelay returns the duration that should be waited before retrying to
// send a query, given the number of retries that have already been made.
//
// The delay starts at 0, since most connection problems can be resolved by
// calling Refresh on the mgo.Session without waiting. It then increases
// quadratically with the number of retries.
func retryDelay(retries int) time.Duration {
	return time.Duration(retries*retries) * time.Second
}

// Collection is an interface over the mgo.Collection type, making possible
// to mock database to handle tests nicely.
// Also, this interface should be used when you want to handle automatic
// retry in case you're in a cluster for example (not yet implemented).
type Collection interface {
	GetByID(id string, out interface{}) error
	Insert(in interface{}) error
	Update(id string, in interface{}) error
	Remove(id string) error
	Upsert(id string, in interface{}) error
	// RemoveAll returns the number of documents removed from collection
	RemoveAll() (int, error)

	// Get is like Find().All()
	Get(filter map[string]interface{}, out interface{}) error

	// GetSorted returns all the documents that match a filter, sorted
	// according to the provided field names. To sort in descending order, add
	// a "-" before the name of the field.
	GetSorted(filter map[string]interface{}, sortFields []string, out interface{}) error

	// GetOne is like Find().One()
	GetOne(filter map[string]interface{}, selector interface{}, out interface{}) error

	// GetFirst returns the first document matching a filter, sorted according
	// to the provided field name.
	GetFirst(filter map[string]interface{}, sortField string, out interface{}) error

	Count(query ...interface{}) (int, error)

	NewBulk(size int) bulk.Bulk

	// Aggregate runs a data-processing pipeline
	// A pipeline calculates aggregate values for the data in a collection
	// Further details : https://docs.mongodb.com/manual/reference/operator/aggregation-pipeline/
	Aggregate(pipeline interface{}, out interface{}) error

	// Pipe to iterate aggregation
	Pipe(pipeline interface{}) *mgo.Pipe
}

type mgoCollection struct {
	coll *mgo.Collection
}

// handleError handles an error returned by a MongoDB query. It logs the error,
// and, if needed, refreshes the connection to the MongoDB database.
//
// It returns false if the query that returned the error should be run again.
// It returns true if there is no error, or if the error should be propagated.
func (c *mgoCollection) handleError(err error, retries int) bool {
	if !shouldTriggerRefresh(err) {
		if retries > 0 {
			log.Printf("reconnection successful")
		}
		return true
	}

	log.Printf(
		"error running MongoDB command (%d/%d retries): (%T) %+v",
		retries, maxRetries, err, err)

	if retries < maxRetries {
		delay := retryDelay(retries)
		log.Printf("waiting %s before retrying", delay)
		time.Sleep(delay)

		c.coll.Database.Session.Refresh()
	}

	return false
}

func (c *mgoCollection) Upsert(id string, in interface{}) error {
	var err error

	for retries := 0; retries <= maxRetries; retries++ {
		_, err = c.coll.UpsertId(id, in)

		if c.handleError(err, retries) {
			break
		}
	}

	return WrapError(err)
}

func (c *mgoCollection) GetByID(id string, out interface{}) error {
	var err error

	for retries := 0; retries <= maxRetries; retries++ {
		err = c.coll.FindId(id).One(out)

		if c.handleError(err, retries) {
			break
		}
	}

	return WrapError(err)
}

func (c *mgoCollection) Insert(in interface{}) error {
	var err error

	for retries := 0; retries <= maxRetries; retries++ {
		err = c.coll.Insert(in)

		if c.handleError(err, retries) {
			break
		}
	}

	return WrapError(err)
}

func (c *mgoCollection) Update(id string, in interface{}) error {
	var err error

	for retries := 0; retries <= maxRetries; retries++ {
		err = c.coll.UpdateId(id, in)

		if c.handleError(err, retries) {
			break
		}
	}

	return WrapError(err)
}

func (c *mgoCollection) Remove(id string) error {
	var err error

	for retries := 0; retries <= maxRetries; retries++ {
		err = c.coll.RemoveId(id)

		if c.handleError(err, retries) {
			break
		}
	}

	return WrapError(err)
}

func (c *mgoCollection) RemoveAll() (int, error) {
	var info *mgo.ChangeInfo
	var err error

	for retries := 0; retries <= maxRetries; retries++ {
		info, err = c.coll.RemoveAll(nil)

		if c.handleError(err, retries) {
			break
		}
	}

	return info.Removed, WrapError(err)
}

func (c *mgoCollection) Get(filter map[string]interface{}, out interface{}) error {
	var err error

	for retries := 0; retries <= maxRetries; retries++ {
		err = c.coll.Find(filter).All(out)

		if c.handleError(err, retries) {
			break
		}
	}

	return WrapError(err)
}

func (c *mgoCollection) GetSorted(filter map[string]interface{}, sortFields []string, out interface{}) error {
	var err error

	for retries := 0; retries <= maxRetries; retries++ {
		err = c.coll.Find(filter).Sort(sortFields...).All(out)

		if c.handleError(err, retries) {
			break
		}
	}

	return WrapError(err)
}

func (c *mgoCollection) GetOne(filter map[string]interface{}, selector interface{}, out interface{}) error {
	var err error

	for retries := 0; retries <= maxRetries; retries++ {
		err = c.coll.Find(filter).Select(selector).One(out)

		if c.handleError(err, retries) {
			break
		}
	}

	return WrapError(err)
}

func (c *mgoCollection) GetFirst(filter map[string]interface{}, sortField string, out interface{}) error {
	var err error

	for retries := 0; retries <= maxRetries; retries++ {
		err = c.coll.Find(filter).Sort(sortField).One(out)

		if c.handleError(err, retries) {
			break
		}
	}

	return WrapError(err)
}

func (c *mgoCollection) NewBulk(size int) bulk.Bulk {
	return bulk.New(c.coll, size)
}

func (c *mgoCollection) Aggregate(pipeline interface{}, out interface{}) error {
	var err error

	pipe := c.coll.Pipe(pipeline)
	for retries := 0; retries <= maxRetries; retries++ {
		err = pipe.All(out)

		if c.handleError(err, retries) {
			break
		}
	}

	return WrapError(err)
}

func (c *mgoCollection) Pipe(pipeline interface{}) *mgo.Pipe {
	return c.coll.Pipe(pipeline)
}

func (c *mgoCollection) Count(query ...interface{}) (int, error) {
	var findQuery interface{}
	if len(query) == 1 {
		findQuery = query[0]
	} else if len(query) > 1 {
		panic("too much arguments")
	}

	var count int
	var err error

	for retries := 0; retries <= maxRetries; retries++ {
		if findQuery != nil {
			count, err = c.coll.Find(findQuery).Count()
		} else {
			count, err = c.coll.Count()
		}

		if c.handleError(err, retries) {
			break
		}
	}

	return count, WrapError(err)
}

func FromMgo(coll *mgo.Collection) Collection {
	return &mgoCollection{
		coll: coll,
	}
}
