# go-deptree

This repo exists to create the simplest go program that will output your dependencies in a tree view.

It will run `go mod graph` and send the output to `go-mod-graph.txt`. The output file will be processed
and rendered as a tree view. The output also excludes dependencies from `golang.org`. There's no way to
change this behavior at this time.

There are a few flags you can pass:
* maxDepth - how deeply to recurse into the deps. defaults to max int
* verbose - just a few extra lines of output
* version - prints the version and exits
* includeVersion - adds the version of the dependency to the tree. defaults to false
* hideSkipReason - suppresses the 'previously seen' and child dependency skip counts

For example, running `go run main.go` on this project renders one dependency on `testify` 
and another dependency on `github.com/json-iterator/go` to demonstrate some dependencies.

```
$ go run main.go
github.com/dovholuknf/go-deptree
     ├─ github.com/json-iterator/go
     │   ├─ github.com/davecgh/go-spew
     │   ├─ github.com/google/gofuzz
     │   ├─ github.com/modern-go/concurrent
     │   ├─ github.com/modern-go/reflect2
     │   └─ github.com/stretchr/testify
     │       ├─ github.com/davecgh/go-spew <previously seen>
     │       ├─ github.com/pmezard/go-difflib
     │       └─ github.com/stretchr/objx
     ├─ github.com/stretchr/testify <previously seen - skipping 4 children>
     └─ go
         └─ toolchain
```

## Installing
To install, execute:
```
GO111MODULE=off go get -u github.com/dovholuknf/go-deptree
go install github.com/dovholuknf/go-deptree@latest
```

## Running
### Defaults
```
$ go-deptree
github.com/dovholuknf/go-deptree
     ├─ github.com/json-iterator/go
     │   ├─ github.com/davecgh/go-spew
     │   ├─ github.com/google/gofuzz
     │   ├─ github.com/modern-go/concurrent
     │   ├─ github.com/modern-go/reflect2
     │   └─ github.com/stretchr/testify
     │       ├─ github.com/davecgh/go-spew <previously seen>
     │       ├─ github.com/pmezard/go-difflib
     │       └─ github.com/stretchr/objx
     ├─ github.com/stretchr/testify <previously seen - skipping 4 children>
     └─ go
         └─ toolchain
```

### Max Depth
```
$ go-deptree -maxDepth=1
Processing with maxDepth: 1
github.com/dovholuknf/go-deptree
     ├─ github.com/json-iterator/go
     ├─ github.com/stretchr/testify
     └─ go
```

### Include Dependency Versions
```
$ go-deptree -includeVersion
github.com/dovholuknf/go-deptree
     ├─ github.com/json-iterator/go@v1.1.12
     │   ├─ github.com/davecgh/go-spew@v1.1.1
     │   ├─ github.com/google/gofuzz@v1.0.0
     │   ├─ github.com/modern-go/concurrent@v0.0.0-20180228061459-e0a39a4cb421
     │   ├─ github.com/modern-go/reflect2@v1.0.2
     │   └─ github.com/stretchr/testify@v1.3.0
     │       ├─ github.com/davecgh/go-spew@v1.1.0
     │       ├─ github.com/pmezard/go-difflib@v1.0.0
     │       └─ github.com/stretchr/objx@v0.1.0
     ├─ github.com/stretchr/testify@v1.8.4
     │   ├─ github.com/davecgh/go-spew@v1.1.1 <previously seen>
     │   ├─ github.com/pmezard/go-difflib@v1.0.0 <previously seen>
     │   ├─ github.com/stretchr/objx@v0.5.0
     │   └─ gopkg.in/yaml.v3@v3.0.1
     │       └─ gopkg.in/check.v1@v0.0.0-20161208181325-20d25e280405
     └─ go@1.21.5
         └─ toolchain@go1.21.5
```

### Include Version with MaxDepth
```
$ go-deptree -maxDepth=1 -includeVersion
Processing with maxDepth: 1
github.com/dovholuknf/go-deptree
     ├─ github.com/json-iterator/go@v1.1.12
     ├─ github.com/stretchr/testify@v1.8.4
     └─ go@1.21.5
```
