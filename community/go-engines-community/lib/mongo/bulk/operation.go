package bulk

// OpInsert handle creation of a bulk insert
type OpInsert interface {
	// Data is the document to be inserted
	Data() interface{}
}

// OpUpdate handle creation of a bulk update
type OpUpdate interface {
	// Selector is a mongo filter used to select documents to be updated.
	Selector() interface{}
	// Data is the query used to perform the update
	Data() interface{}
}

type opInsert struct {
	data interface{}
}

func (o opInsert) Data() interface{} {
	return o.data
}

type opUpdate struct {
	selector interface{}
	data     interface{}
}

func (o opUpdate) Data() interface{} {
	return o.data
}

func (o opUpdate) Selector() interface{} {
	return o.selector
}

// NewInsert creates an insert operation.
func NewInsert(data interface{}) OpInsert {
	return opInsert{
		data: data,
	}
}

// NewUpdate creates an update operation.
//
// selector is a MongoDB filter to match document(s) to be updated with
// the data parameter.
func NewUpdate(selector interface{}, data interface{}) OpUpdate {
	return opUpdate{
		selector: selector,
		data:     data,
	}
}
