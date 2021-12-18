package lockmanager

import (
	"log"
	"sync"
	"sync/atomic"
)

type lmData struct {
	wg    sync.WaitGroup
	count int32
}
type LockManager struct {
	locks map[string]*lmData
	mu    sync.Mutex
}

func New() *LockManager {
	lm := LockManager{
		locks: make(map[string]*lmData),
	}
	return &lm
}

func (l *LockManager) getLock(key string) *lmData {
	l.mu.Lock()
	defer l.mu.Unlock()

	ld, ok := l.locks[key]
	if !ok {
		ld = &lmData{wg: sync.WaitGroup{}}
		l.locks[key] = ld
	}
	return ld
}

func (l *LockManager) cleanup(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	ld, ok := l.locks[key]
	if ok && atomic.LoadInt32(&ld.count) == 0 {
		delete(l.locks, key)
	}
}

// Do : lock on filename and run dofunc()
func (l *LockManager) Do(filename string, dofunc func()) {
	ld := l.getLock(filename)

	atomic.AddInt32(&ld.count, 1)
	log.Println("wait", filename)
	ld.wg.Wait() // wait for existing to finish
	ld.wg.Add(1) // block the rest
	log.Println("doing", filename)
	dofunc()
	ld.wg.Done()
	atomic.AddInt32(&ld.count, -1)
	l.cleanup(filename)
}
