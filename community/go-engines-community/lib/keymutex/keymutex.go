// Package keymutex contains mutex that allows locking by key.
package keymutex

import (
	"fmt"
	"sync"
)

type KeyMutex interface {
	Lock(string)
	Unlock(string) error
	LockMultiple(keys ...string)
	UnlockMultiple(keys ...string) error
}

func New() KeyMutex {
	return &keyMutex{
		locks: make(map[string]*mutexWithCount),
	}
}

type keyMutex struct {
	locksMx sync.Mutex
	locks   map[string]*mutexWithCount
}

type mutexWithCount struct {
	count int64
	mx    sync.Mutex
}

func (m *keyMutex) Lock(key string) {
	mx := m.getOrCreateMx(key)
	mx.Lock()
}

func (m *keyMutex) Unlock(key string) error {
	mx, err := m.getAndDeleteMx(key)
	if err != nil {
		return err
	}

	mx.Unlock()
	return nil
}

func (m *keyMutex) LockMultiple(keys ...string) {
	mutexes := m.getOrCreateMutexes(keys)

	for _, mx := range mutexes {
		mx.Lock()
	}
}

func (m *keyMutex) UnlockMultiple(keys ...string) error {
	mutexes, err := m.getAndDeleteMutexes(keys)
	if err != nil {
		return err
	}

	for _, mx := range mutexes {
		mx.Unlock()
	}

	return nil
}

func (m *keyMutex) getOrCreateMx(key string) *sync.Mutex {
	m.locksMx.Lock()
	defer m.locksMx.Unlock()

	var lock *mutexWithCount
	var ok bool
	if lock, ok = m.locks[key]; !ok {
		lock = &mutexWithCount{}
		m.locks[key] = lock
	}

	lock.count++

	return &lock.mx
}

func (m *keyMutex) getAndDeleteMx(key string) (*sync.Mutex, error) {
	m.locksMx.Lock()
	defer m.locksMx.Unlock()

	if lock, ok := m.locks[key]; ok {
		lock.count--

		if lock.count == 0 {
			delete(m.locks, key)
		}

		return &lock.mx, nil
	}

	return nil, fmt.Errorf("mutex for %q not found", key)
}

func (m *keyMutex) getOrCreateMutexes(keys []string) []*sync.Mutex {
	m.locksMx.Lock()
	defer m.locksMx.Unlock()

	mutexes := make([]*sync.Mutex, len(keys))

	for i, key := range keys {
		var lock *mutexWithCount
		var ok bool

		if lock, ok = m.locks[key]; !ok {
			lock = &mutexWithCount{}
			m.locks[key] = lock
		}

		lock.count++
		mutexes[i] = &lock.mx
	}

	return mutexes
}

func (m *keyMutex) getAndDeleteMutexes(keys []string) ([]*sync.Mutex, error) {
	m.locksMx.Lock()
	defer m.locksMx.Unlock()

	for _, key := range keys {
		if _, ok := m.locks[key]; !ok {
			return nil, fmt.Errorf("mutex for %q not found", key)
		}
	}

	mutexes := make([]*sync.Mutex, len(keys))

	for i, key := range keys {
		lock := m.locks[key]
		lock.count--

		if lock.count == 0 {
			delete(m.locks, key)
		}

		mutexes[i] = &lock.mx
	}

	return mutexes, nil
}
