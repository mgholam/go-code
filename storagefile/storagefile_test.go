package storagefile_test

import (
	"bb/storagefile"
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type Book struct {
	ID     int       `json:"id"`
	Title  string    `json:"name"`
	Author string    `json:"author"`
	Rating int       `json:"rating"`
	Date   time.Time `json:"date" gorm:"column:date"`
}

func Test_run_100k_read_write(t *testing.T) {

	sf, e := storagefile.Open("docs.dat")
	if e != nil {
		panic(e)
	}
	defer sf.Close()
	defer func() {
		os.Remove("docs.dat")
		os.Remove("docs.dat.idx")
	}()

	count := 100_000

	fmt.Println("saving count", count)
	dosave(sf, count)
	count = int(sf.Count())
	fmt.Println("reading count", count)
	doread(sf, count)

	h, _ := sf.GetHeader(10)
	fmt.Println(h)

	tt, ss, _ := sf.GetString(10)
	fmt.Println(tt, ss)

	tt, b, _ := sf.Get(10)
	fmt.Println(tt, b)
}

func Test_rebuild(t *testing.T) {
	defer func() {
		os.Remove("rebuild.dat")
		os.Remove("rebuild.dat.idx")
	}()
	os.Remove("rebuild.dat")
	sf, e := storagefile.Open("rebuild.dat")
	if e != nil {
		panic(e)
	}

	for i := 1; i <= 1000; i++ {
		sf.Save("22", []byte("111111"))
	}
	for i := 1; i <= 1000; i++ {
		tt, s, _ := sf.GetString(int64(i))
		if s != "111111" || tt != "22" {
			t.Fail()
		}
	}

	sf.Close()
	os.WriteFile("rebuild.dat.dirty", []byte("hello"), 0644)
	sf, e = storagefile.Open("rebuild.dat")
	if e != nil {
		panic(e)
	}
	for i := 1; i <= 1000; i++ {
		tt, s, e := sf.GetString(int64(i))
		if e != nil {
			t.Log(e)
			t.Fail()
		}
		if s != "111111" || tt != "22" {
			t.Fail()
		}
	}

	sf.Close()

}

func Test_open_close_twice(t *testing.T) {
	defer func() {
		os.Remove("oc2.dat")
		os.Remove("oc2.dat.idx")
	}()
	os.Remove("oc2.dat")
	sf, e := storagefile.Open("oc2.dat")
	if e != nil {
		panic(e)
	}

	sf.Save("11", []byte("1111111"))
	sf.Save("11", []byte("1111111"))
	sf.Close()

	sf, e = storagefile.Open("oc2.dat")
	if e != nil {
		panic(e)
	}

	sf.Save("22", []byte("2222222"))
	sf.Save("22", []byte("2222222"))

	if sf.Count() != int64(4) {
		t.Fail()
	}

	for i := 1; i <= 4; i++ {
		h, e := sf.GetHeader(int64(i))
		if e != nil {
			t.Log(e)
			t.Fail()
		}
		if h.Id != int64(i) {
			t.Log("id mismatch", i, h.Id)
			t.Fail()
		}
	}
	tt, d, e := sf.GetString(1)
	if e != nil {
		t.Log(e)
		t.Fail()
	}
	if tt != "11" || d != "1111111" {
		t.Log("type and data mismatch")
		t.Fail()
	}
	sf.Close()
}

func Test_invalid(t *testing.T) {
	defer func() {
		os.Remove("inv.dat")
		os.Remove("inv.dat.idx")
	}()
	os.Remove("inv.dat")
	sf, e := storagefile.Open("inv.dat")
	if e != nil {
		panic(e)
	}
	defer sf.Close()

	_, _, e = sf.Get(1)
	if e != nil {
		t.Log(e)
	}

	_, _, e = sf.GetString(1)
	if e != nil {
		t.Log(e)
	}
}

func doread(sf *storagefile.StorageFile, count int) {
	t := time.Now()
	for i := 1; i <= count; i++ {
		h, e := sf.GetHeader(int64(i))
		if e != nil {
			fmt.Println("err", e, i)
			continue
		}
		if h.Type != "/api/book" {
			// if h.Id%100_000 != int64(i) {
			fmt.Println("data not matching", i)
		}
	}
	fmt.Println("read time =", time.Since(t))
}

func dosave(sf *storagefile.StorageFile, count int) {

	t := time.Now()

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	// json.NewDecoder()
	// jsoniter.ConfigFastest

	for i := 1; i <= count; i++ {
		book := Book{
			ID:     i,
			Author: "tolkien",
			Title:  "lord of the rings",
			Rating: 5,
			Date:   now(),
		} // new(Book)
		b, _ := json.Marshal(book)
		// b, _ := json.Marshal(book)
		// b, _ := easyjson.Marshal(book)

		sf.Save("/api/book", b)
	}
	fmt.Println("save time =", time.Since(t))
}

func now() time.Time {
	var tv syscall.Timeval
	syscall.Gettimeofday(&tv)
	return time.Unix(0, syscall.TimevalToNsec(tv))
}