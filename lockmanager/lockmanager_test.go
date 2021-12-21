package lockmanager_test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/mgholam/go-code/lockmanager"
)

func Test(t *testing.T) {

	var i int32
	var wg = sync.WaitGroup{}
	wg.Add(4)

	go lockmanager.Do("file2", func() {
		// fmt.Println(lm.Locks)
		time.Sleep(2 * time.Second)
		atomic.AddInt32(&i, 1)
		t.Log("1 done")
		wg.Done()
	})
	go lockmanager.Do("file1", func() {
		// fmt.Println(lm.Locks)

		time.Sleep(2 * time.Second)

		atomic.AddInt32(&i, 1)
		t.Log("2 done")
		wg.Done()
	})
	go lockmanager.Do("file2", func() {
		// fmt.Println(lm.Locks)
		// fmt.Println("3 locking file")
		time.Sleep(2 * time.Second)
		atomic.AddInt32(&i, 1)

		t.Log("3 done")
		wg.Done()
	})
	go lockmanager.Do("file1", func() {
		// fmt.Println(lm.Locks)
		fmt.Println("4 locking file")
		time.Sleep(2 * time.Second)
		atomic.AddInt32(&i, 1)
		t.Log("4 done")
		wg.Done()
	})

	wg.Wait()
	t.Log(i)
	// fmt.Println(lm.locks)
}
