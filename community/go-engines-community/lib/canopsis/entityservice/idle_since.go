package entityservice

import "sync"

type ServicesIdleSinceMap struct {
	idleMap map[string]int64
	mutex   sync.Mutex
}

func NewServicesIdleSinceMap() ServicesIdleSinceMap {
	return ServicesIdleSinceMap{
		idleMap: make(map[string]int64),
		mutex:   sync.Mutex{},
	}
}

func (m *ServicesIdleSinceMap) Mark(id string, newIdleSince int64) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	oldIdleSince := m.idleMap[id]
	if oldIdleSince == 0 || oldIdleSince > newIdleSince {
		m.idleMap[id] = newIdleSince
		return true
	}

	return false
}
