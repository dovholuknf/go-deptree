package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func TestName(t *testing.T) { //function just exists to add a second dep with testify as a dependency
	val := []byte(`{"ID":1,"Name":"Reds","Colors":["Crimson","Red","Ruby","Maroon"]}`)
	jsoniter.Get(val, "Colors", 0).ToString()
}

func setup(hideSkipReason, includeVersion, verbose, version bool) (*DepTreeOpts, *strings.Builder) {
	opts := NewDepTreeOpts()
	opts.hideSkipReason = hideSkipReason
	opts.includeVersion = includeVersion
	opts.verbose = verbose
	opts.version = version
	sb := &strings.Builder{}
	opts.writer = sb
	return opts, sb
}

func TestDefaults(t *testing.T) {
	opts, result := setup(false, false, false, false)
	process(opts)
	output := result.String()
	fmt.Println(output)

	assert.Contains(t, output, "github.com/davecgh/go-spew <previously seen>")
	assert.Contains(t, output, "github.com/stretchr/testify <previously seen - skipping 4 children>")
}

func TestVerbose(t *testing.T) {
	opts, result := setup(false, false, true, false)
	process(opts)
	output := result.String()
	fmt.Println(output)

	assert.Contains(t, output, "Reading go mod graph file")
}

func TestVersion(t *testing.T) {
	opts, result := setup(false, false, false, true)
	process(opts)
	output := result.String()
	fmt.Println(output)

	assert.Contains(t, output, ver)
}

func TestHideSkipReason(t *testing.T) {
	opts, result := setup(true, false, false, false)
	process(opts)
	output := result.String()
	fmt.Println(output)

	assert.NotContains(t, output, "previously seen")
}

func TestMaxDepth(t *testing.T) {
	opts, result := setup(false, false, false, false)
	opts.depth = 1
	process(opts)
	output := result.String()
	fmt.Println(output)

	assert.Contains(t, output, "github.com/stretchr/testify")
	assert.NotContains(t, output, "github.com/davecgh/go-spew")
}
