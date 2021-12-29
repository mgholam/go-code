package zipfile_test

import (
	"os"
	"strings"
	"testing"

	"github.com/mgholam/go-code/zipfile"
)

func Test(t *testing.T) {
	os.Remove("ziptest.zip")
	os.Remove("zz.zip")

	// zipfile.UniqueFiles = true
	zipfile.AddFile("ziptest.zip", "zipfile.go")
	// zipfile.AddFile("ziptest.zip", "../genthumb/readme.md")
	zipfile.AddFile("ziptest.zip", "zipfile.go")

	zipfile.Add("zz.zip", "aa.txt", strings.NewReader("aa"))
	zipfile.Add("zz.zip", "aa.txt", strings.NewReader("bb"))
	zipfile.Add("zz.zip", "aa.txt", strings.NewReader("cc"))
}
