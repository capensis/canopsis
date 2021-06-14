package context

import (
	"sync"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
)

type safeEC struct {
	unsafe EnrichmentCenter
	l      sync.RWMutex
}

func (s *safeEC) Flush() error {
	s.l.Lock()
	defer s.l.Unlock()
	return s.unsafe.Flush()
}

func (s *safeEC) LoadWatchers() error {
	s.l.Lock()
	defer s.l.Unlock()
	return s.unsafe.LoadWatchers()
}

func (s *safeEC) Handle(event types.Event, ef EnrichFields) *types.Entity {
	s.l.Lock()
	defer s.l.Unlock()
	return s.unsafe.Handle(event, ef)
}

func (s *safeEC) Get(event types.Event) (*types.Entity, error) {
	s.l.Lock()
	defer s.l.Unlock()
	return s.unsafe.Get(event)
}

func (s *safeEC) Update(entity types.Entity) types.Entity {
	s.l.Lock()
	defer s.l.Unlock()
	return s.unsafe.Update(entity)
}

func (s *safeEC) EnrichResourceInfoWithComponentInfo(event *types.Event, entity *types.Entity) error {
	s.l.Lock()
	defer s.l.Unlock()
	return s.unsafe.EnrichResourceInfoWithComponentInfo(event, entity)
}

// NewSafeEnrichmentCenter wrap EnrichmentCenter with a sync.RWMutex
func NewSafeEnrichmentCenter(ec EnrichmentCenter) EnrichmentCenter {
	sec := safeEC{
		unsafe: ec,
		l:      sync.RWMutex{},
	}

	return &sec
}
