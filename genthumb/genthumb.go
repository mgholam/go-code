//usr/bin/go run $0 $@ ; exit

package main

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var donecount int = 0

// var run bool = true;

func main() {
	args := os.Args
	fn, _ := os.UserHomeDir()
	var doublecmdThumbDir string = ".cache/doublecmd/thumbnails"
	var dimension int16 = 150
	doublecmdThumbDir = path.Join(fn, doublecmdThumbDir)

	if len(args) < 3 {
		fmt.Println("\nUsage : genthumbs <dimension> <picture directory>")
		return
	}

	fmt.Println("thumb folder : ", doublecmdThumbDir)
	dim, e := strconv.Atoi(args[1])
	if e != nil {
		fmt.Println(e)
		return
	}
	dimension = int16(dim)

	count := 0
	dir, _ := filepath.Abs(args[2])
	var wg = sync.WaitGroup{}
	start := time.Now()

	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			genthumbs(path, doublecmdThumbDir, dimension)
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

func genthumbs(path string, doublecmdThumbDir string, dimension int16) {
	f, e := os.Stat(path)
	if e != nil || f.IsDir() {
		return
	}

	excludelist := []string{".mp4", ".gif"}

	for _, s := range excludelist {
		if strings.HasSuffix(path, s) {
			return
		}
	}

	b := md5.Sum([]byte(path))
	tfn := filepath.Join(doublecmdThumbDir, fmt.Sprintf("%x.jpg", b))

	if isExist(tfn) {
		return // found
	}

	fmt.Println("processing : ", path)

	// do stuff
	e = exec.Command("convert", "-thumbnail", fmt.Sprintf("%dx%d", dimension, dimension), "-strip", "-quality", "80", path, tfn).Run()
	if e != nil {
		fmt.Println(e)
		return
	}
	writeMetaData(path, tfn, dimension)
}

func writeMetaData(path string, fn string, dimension int16) {
	st, e := os.Stat(path)
	if e != nil {
		fmt.Println("error : ", e)
		return
	}

	f, e := os.OpenFile(fn, os.O_RDWR, 0644)
	if e != nil {
		fmt.Println("error : ", e)
		return
	}
	defer f.Close()

	iend, _ := f.Seek(0, os.SEEK_END)

	f.Write([]byte{0, 0})
	f.WriteString("#THUMB")
	s := "file://" + path
        s = strings.ReplaceAll(s, " ", "%20") // kludge
	var l int32 = int32(len(s))
	binary.Write(f, binary.LittleEndian, l) // s len
	f.WriteString(s)                        // to uri

	binary.Write(f, binary.LittleEndian, st.Size())                    // org file len
	binary.Write(f, binary.LittleEndian, encodedatetime(st.ModTime())) // org file mod time
	binary.Write(f, binary.LittleEndian, dimension)                    // x
	binary.Write(f, binary.LittleEndian, dimension)                    // y
	binary.Write(f, binary.LittleEndian, int32(iend))                  //iend

	donecount++
}

func isExist(fn string) bool {
	_, e := os.Stat(fn)
	return e == nil
}

func encodedatetime(mt time.Time) float64 {
	return encodedate(int32(mt.Year()), int32(mt.Month()), int32(mt.Day())) +
		encodetime(int32(mt.Hour()), int32(mt.Minute()), int32(mt.Second()), int32(mt.Nanosecond())/1e6)
}

func encodedate(y int32, m int32, d int32) float64 {
	var date float64 = 0

	if m > 2 {
		m -= 3
	} else {
		m += 9
		y--
	}
	c := int64(y / 100)
	ya := int64(y) - 100*c
	dt := (146097*c)/4 + (1461*ya)/4 + (153*int64(m)+2)/5 + int64(d)
	dt -= 693900
	date = float64(dt)

	return date
}

func encodetime(h int32, m int32, s int32, ms int32) float64 {

	i := float64(h)*3600000 + float64(m)*60000 + float64(s)*1000 + float64(ms)
	i = i / (24 * 3600 * 1000)
	return i
}

// func conImage() {
// imagePath, _ := os.Open(path)
// defer imagePath.Close()
// srcImage, _, e := image.Decode(imagePath)
// if e != nil {
// 	fmt.Println(e)
// 	return
// }

// dstImage := image.NewRGBA(image.Rect(0, 0, dimension, dimension))
// // Thumbnail function of Graphics
// graphics.Thumbnail(dstImage, srcImage)

// newImage, _ := os.Create(tfn)
// defer newImage.Close()
// // jpeg.Encode(newImage, dstImage, &jpeg.Options{jpeg.DefaultQuality})
// }

/*

	below original delphi code

*/

// Function TryEncodeDateTime(const AYear, AMonth, ADay, AHour, AMinute, ASecond, AMilliSecond: Word; out AValue: TDateTime): Boolean;
// var
//  tmp : TDateTime;
// begin
//   Result:=TryEncodeDate(AYear,AMonth,ADay,AValue);
//   Result:=Result and TryEncodeTime(AHour,AMinute,ASecond,Amillisecond,Tmp);
//   If Result then
//     Avalue:=ComposeDateTime(AValue,Tmp);
// end;

// Function TryEncodeDate(Year,Month,Day : Word; Out Date : TDateTime) : Boolean;
// var
//   c, ya: cardinal;
// begin
//   Result:=(Year>0) and (Year<10000) and
//           (Month in [1..12]) and
//           (Day>0) and (Day<=MonthDays[IsleapYear(Year),Month]);
//  If Result then
//    begin
//      if month > 2 then
//       Dec(Month,3)
//      else
//       begin
//         Inc(Month,9);
//         Dec(Year);
//       end;
//      c:= Year DIV 100;
//      ya:= Year - 100*c;
//      Date := (146097*c) SHR 2 + (1461*ya) SHR 2 + (153*cardinal(Month)+2) DIV 5 + cardinal(Day);
//      // Note that this line can't be part of the line above, since TDateTime is
//      // signed and c and ya are not
//      Date := Date - 693900;
//    end
// end;

// function TryEncodeTime(Hour, Min, Sec, MSec:word; Out Time : TDateTime) : boolean;
// begin
//   Result:=(Hour<24) and (Min<60) and (Sec<60) and (MSec<1000);
//   If Result then
//     Time:=TDateTime(cardinal(Hour)*3600000+cardinal(Min)*60000+cardinal(Sec)*1000+MSec)/MSecsPerDay;
// end;
