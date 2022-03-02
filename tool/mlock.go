package tool

import "sync"

var m *mlock

type mlock struct {
	lock  sync.RWMutex
	locks map[string]*sync.RWMutex
}

func initmlock() {
	m = new(mlock)
	m.locks = make(map[string]*sync.RWMutex)
}

func init() {
	initmlock()
}

func getLock(key string) *sync.RWMutex {
	m.lock.RLock()
	lock, ok := m.locks[key]
	m.lock.RUnlock()
	if ok {
		return lock
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	lock, ok = m.locks[key]
	if ok {
		return lock
	}
	m.locks[key] = new(sync.RWMutex)
	return m.locks[key]
}

func Lock(key string) {
	getLock(key).Lock()
}

func Unlock(key string) {
	getLock(key).Unlock()
}

func RLock(key string) {
	getLock(key).RLock()
}

func RUnlock(key string) {
	getLock(key).RUnlock()
}

func Use(key string, f func()) {
	Lock(key)
	defer Unlock(key)
	f()
}

func RUse(key string, f func()) {
	RLock(key)
	defer RUnlock(key)
	f()
}
