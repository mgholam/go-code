# zipfile

Missing features from the standard library `zip` package:

- add files to an existing zip file
- adds the file time as a comment since the original zip code uses MSDOS format dates and times and only saves even seconds.
- `zipfile` does not overwrite existing files in the zip file, so you can have backups of file changes within the zip file.

```go
// add a file on disk
zipfile.AddFile("ziptest.zip", "zipfile.go")

// use an io.Reader
zipfile.Add("zz.zip", "aa.txt", strings.NewReader("aa"))
```

