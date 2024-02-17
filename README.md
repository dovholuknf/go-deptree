# go-deptree

This repo exists to create the simplest go program that will output your dependencies in a tree view.

It will run `go mod graph` and send the output to `go-mod-graph.txt`. The output file will be processed
and rendered as a tree view. The output also excludes dependencies from `golang.org`. There's no way to
change this behavior at this time.

There are three simple flags you can pass:
* maxDepth - how deeply to recurse into the deps. defaults to max int
* verbose - just a few extra lines of output
* version - prints the version and exits
* includeVersion - adds the version of the dependency to the tree. defaults to false

For example, running `go run main.go` on this project renders no dependencies other than go itself:

````
$ go run main.go  -verbose -maxDepth=100
Processing with maxDepth: 100
Output of go mod graph written to: ./go-mod-graph.txt

github.com/dovholuknf/go-deptree
    └── go@1.21.5
        └── toolchain@go1.21.5
```

## Installing
To install/update execute:
```
GO111MODULE=off go get -u github.com/dovholuknf/go-deptree
go install github.com/dovholuknf/go-deptree@latest
go-deptree -version
```
