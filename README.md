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

For example, running `go run main.go` on this project renders one dependency on testify 
just to see it in action (no tests exist yet):

```
$ go run main.go
github.com/dovholuknf/go-deptree
     ├─ github.com/stretchr/testify
     │   ├─ github.com/davecgh/go-spew
     │   ├─ github.com/pmezard/go-difflib
     │   ├─ github.com/stretchr/objx
     │   └─ gopkg.in/yaml.v3
     │       └─ gopkg.in/check.v1
     └─ go
         └─ toolchain
```

## Installing
To install execute:
```
GO111MODULE=off go get -u github.com/dovholuknf/go-deptree
go install github.com/dovholuknf/go-deptree@latest
```

## Running
### Defaults
```
$ go-deptree
github.com/dovholuknf/go-deptree
     ├─ github.com/stretchr/testify
     │   ├─ github.com/davecgh/go-spew
     │   ├─ github.com/pmezard/go-difflib
     │   ├─ github.com/stretchr/objx
     │   └─ gopkg.in/yaml.v3
     │       └─ gopkg.in/check.v1
     └─ go
         └─ toolchain
```

### Max Depth
```
$ go-deptree -maxDepth=1
Processing with maxDepth: 1
github.com/dovholuknf/go-deptree
     ├─ github.com/stretchr/testify
     └─ go
```

### Include Version
```
$ go-deptree -maxDepth=1 -includeVersion
Processing with maxDepth: 1
github.com/dovholuknf/go-deptree
     ├─ github.com/stretchr/testify@v1.8.4
     └─ go@1.21.5
```