# genthumb.go

A utility to generate `double commander` thumbnails from a directory of images

```sh
# shebang autorun mode, generate 150x150 thumbnails from /wallpapers folder
$ genthumb.go 150 /wallpapers
```

The code uses `convert` to create thumbnails `sudo <packagemanger> install imagemagick`

# cleanupthumb.go

A utility to check the `double commander` thumbnail cache directory on linux systems for invalid thumbnails and delete them if the original is not present.

```sh
# the script has a shebang header and will autorun
$ cleanupthumb.go
# or normal
$ go run cleanupthumb.go
```

# Interesting parts

- The code uses a `filepath.walker` and starts a `go routine` per file, it also waits every 200 files to finish so you don't get too many files open errors.
- `double commander` adds meta data to the end of the thumbnail for tracking, so a simple thumbnail does not work
- `delphi` uses weird binary formated date/time data so the pascal code was reverse engineered to `go`