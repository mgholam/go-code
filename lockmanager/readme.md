# lockmanager

`lockmanager` ensures serialized execution of a unit of work on a resource identifier. It can be used for :

- blocking usage of a file

```go
var lm = lockmanager.New()

go lm.Do("file1", func() {
		time.Sleep(2 * time.Second)
		fmt.Println("1 done")
})

go lm.Do("file1", func() {
		time.Sleep(2 * time.Second)
		fmt.Println("2 done")
})
```

The string `file1` does not have to be an actual file its just a place holder for synchronization. You can use the code to coordinate using files in a multi thread/go routine application and queued go routines will wait and block until allowed to run.

## Use case

I needed this code for a web server that uses zip files as storage of input data.

## Thread safety

The code is go routine and race safe and uses locks on a `map[]` structure and `sync.WaitGroup` for function execution coordination. The `go run -race` command does not show any errors.

