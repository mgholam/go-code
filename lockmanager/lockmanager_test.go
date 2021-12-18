package lockmanager_test

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/mgholam/go-code/lockmanager"
)

func Test(t *testing.T) {

	var i int32
	var lm = lockmanager.New()


	go lm.Do("file2", func() {
		// fmt.Println(lm.Locks)
		time.Sleep(2 * time.Second)
		atomic.AddInt32(&i, 1)
		fmt.Println("1 done")
	})
	go lm.Do("file1", func() {
		// fmt.Println(lm.Locks)

		time.Sleep(2 * time.Second)

		atomic.AddInt32(&i, 1)
		fmt.Println("2 done")
	})
	go lm.Do("file2", func() {
		// fmt.Println(lm.Locks)
		// fmt.Println("3 locking file")
		time.Sleep(2 * time.Second)
		atomic.AddInt32(&i, 1)

		fmt.Println("3 done")
	})
	go lm.Do("file1", func() {
		// fmt.Println(lm.Locks)
		fmt.Println("4 locking file")
		time.Sleep(2 * time.Second)
		atomic.AddInt32(&i, 1)
		fmt.Println("4 done")
	})

	fmt.Println("Press enter...")
	fmt.Scanln()
	fmt.Println(i)
	// fmt.Println(lm.locks)
}
