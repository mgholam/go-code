//usr/bin/go run $0 $@ ; exit
package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var donecount int = 0

func main() {
	fn, _ := os.UserHomeDir()
	var doublecmdThumbDir string = ".cache/doublecmd/thumbnails"
	doublecmdThumbDir = path.Join(fn, doublecmdThumbDir)

	fmt.Println("thumb folder : ", doublecmdThumbDir)
	count := 0
	var wg = sync.WaitGroup{}
	start := time.Now()

	filepath.Walk(doublecmdThumbDir, func(path string, f os.FileInfo, err error) error {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			checkthumb(path)
			wg.Done()
		}(&wg)
		count++
		if count%200 == 0 { // fix for too many open files
			wg.Wait()
		}
		return nil
	})

	wg.Wait()

	fmt.Println()
	fmt.Println("Time taken : ", time.Since(start))
	fmt.Println("Total : ", count, "files")
	fmt.Println("Processed : ", donecount, "files")
}

func checkthumb(path string) {
	fn, e := os.Stat(path)
	if e != nil || fn.IsDir() {
		return
	}
	tlen := fn.Size()

	f, e := os.Open(path)
	if e != nil {
		fmt.Println("error : ", e)
		return
	}

	defer f.Close()
	f.Seek(-4, os.SEEK_END)
	b := make([]byte, 4)

	var pos int32 = -1
	binary.Read(f, binary.LittleEndian, &pos)
	if pos < 0 {
		fmt.Println("!pos", path)
		return
	}
	len := tlen - int64(pos)
	if len > 2000 || len < 0 {
		fmt.Println("block len error", path)
		return
	}
	b = make([]byte, len)

	f.Seek(int64(pos), os.SEEK_SET)
	f.Read(b)

	hdr := []byte{0, 0, '#', 'T', 'H', 'U', 'M', 'B'}

	for i, c := range hdr {
		if b[i] != c {
			fmt.Println("header err", path)
			return
		}
	}

	slen := int32(binary.LittleEndian.Uint32(b[8:12]))
	if slen < 0 || slen > 3000 {
		fmt.Println("filename len error")
		return
	}
	slen += 12
	filename := string(b[12:slen])
	filename = strings.ReplaceAll(filename, "file://", "")
	filename = strings.ReplaceAll(filename, "%20", " ")

	if isExist(filename) == false {
		fmt.Printf("not exist %x %s %s\r\n", slen, filename, path)

		// delete thumb file
		os.Remove(path)
		return
	}

	donecount++
}

func isExist(fn string) bool {
	_, e := os.Stat(fn)
	return e == nil
}
