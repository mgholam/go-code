package lockmanager

import (
	"log"
	"sync"
	"sync/atomic"
)

var l *lockManager

type lmData struct {
	wg    sync.WaitGroup
	count int32
}
type lockManager struct {
	locks map[string]*lmData
	mu    sync.Mutex
}

func getLock(key string) *lmData {
	l.mu.Lock()
	defer l.mu.Unlock()

	ld, ok := l.locks[key]
	if !ok {
		ld = &lmData{wg: sync.WaitGroup{}}
		l.locks[key] = ld
	}
	return ld
}

func cleanup(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	ld, ok := l.locks[key]
	if ok && atomic.LoadInt32(&ld.count) == 0 {
		delete(l.locks, key)
	}
}

// Do : lock on filename and run dofunc()
func Do(filename string, dofunc func()) {
	if l == nil {
		l = &lockManager{
			locks: make(map[string]*lmData),
		}
	}
	ld := getLock(filename)

	atomic.AddInt32(&ld.count, 1)
	log.Println("wait", filename)
	ld.wg.Wait() // wait for existing to finish
	ld.wg.Add(1) // block the rest
	log.Println("doing", filename)
	dofunc()
	ld.wg.Done()
	atomic.AddInt32(&ld.count, -1)
	cleanup(filename)
}
