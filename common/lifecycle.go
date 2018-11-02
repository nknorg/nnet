package common

import "sync"

// LifeCycle is an abstract type that has thread-safe initialization and shutdown
type LifeCycle struct {
	StartOnce sync.Once
	StopOnce  sync.Once
	stopped   bool
	stopLock  sync.RWMutex
	ready     bool
	readyLock sync.RWMutex
}

// Stop set lifecycle status to stopped
func (l *LifeCycle) Stop() {
	l.stopLock.Lock()
	l.stopped = true
	l.stopLock.Unlock()
}

// IsStopped returns if lifecycle is stopped
func (l *LifeCycle) IsStopped() bool {
	l.stopLock.RLock()
	defer l.stopLock.RUnlock()
	return l.stopped
}

// SetReady sets the ready state of lifecycle
func (l *LifeCycle) SetReady(ready bool) {
	l.readyLock.Lock()
	l.ready = ready
	l.readyLock.Unlock()
}

// IsReady returns if lifecycle is ready
func (l *LifeCycle) IsReady() bool {
	l.readyLock.RLock()
	defer l.readyLock.RUnlock()
	return l.ready
}
