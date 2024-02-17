# go-deptree

This repo exists to create the simplest go program that will output your dependencies in a tree view.

It will run `go mod graph` and send the output to `go-mod-graph.txt`. The output file will be processed
and rendered as a tree view.

For example, running `go run main.go` on this project renders no dependencies other than go itself:

```
$ go run main.go
Command successfully executed. Output written to ./go-mod-graph.txt
github.com/dovholuknf/deptree
    └── go@1.21.5
        └── toolchain@go1.21.5
```

## Installing
Installation is easy. Just run `go install github.com/dovholuknf/go-deptree@latest`