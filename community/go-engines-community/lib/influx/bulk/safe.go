package bulk

import "sync"

type safe struct {
	l      sync.RWMutex
	unsafe Bulk
}

func (b *safe) AddPoints(po ...PointOp) error {
	b.l.Lock()
	defer b.l.Unlock()
	return b.unsafe.AddPoints(po...)
}

func (b *safe) Perform() error {
	b.l.Lock()
	defer b.l.Unlock()
	return b.unsafe.Perform()
}

// NewSafe wraps un unsafe Bulk with a sync.RWMutex
func NewSafe(unsafe Bulk) Bulk {
	return &safe{
		unsafe: unsafe,
	}
}
