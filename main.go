package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

const ver = "1.0.9"

type DepTreeOpts struct {
	version        bool
	verbose        bool
	includeVersion bool
	hideSkipReason bool
	rendered       map[string]string
	goModDeps      map[string]bool
	goModFile      string
	goModGraphFile string
	writer         io.Writer
	depth          int
}

type Node struct {
	Value          string
	Children       []*Node
	Parent         string
	Indirect       bool
	includeVersion bool
}

func (n *Node) Val() string {
	if n.includeVersion {
		return n.Value
	} else {
		return strings.Split(n.Value, "@")[0]
	}
}

func NewDepTreeOpts() *DepTreeOpts {
	return &DepTreeOpts{
		rendered:       make(map[string]string),
		goModDeps:      make(map[string]bool),
		goModFile:      "./go.mod",
		goModGraphFile: "./go-mod-graph.txt",
		writer:         os.Stdout,
		depth:          math.MaxInt,
	}
}

func main() {
	time.Sleep(time.Second)
	opts := NewDepTreeOpts()
	maxDepth := flag.Int("maxDepth", math.MaxInt, "Maximum depth for processing")
	flag.BoolVar(&opts.version, "version", false, "Print the version and exits")
	flag.BoolVar(&opts.verbose, "verbose", false, "Print additional debug-type output")
	flag.BoolVar(&opts.includeVersion, "includeVersion", false, "Adds the version of the dependency to the output")
	flag.BoolVar(&opts.hideSkipReason, "hideSkipReason", false, "Suppresses the 'previously seen' and child dependency skip counts")
	flag.Parse()
	opts.depth = *maxDepth
	process(opts)
}

func process(opts *DepTreeOpts) {
	if opts.version {
		_, _ = fmt.Fprintf(opts.writer, ver)
		return
	}

	if opts.depth < 1 {
		_, _ = fmt.Fprintf(opts.writer, "maxDepth cannot be < 1, using 1 for maxDepth")
		opts.depth = 1
	}
	if opts.depth < math.MaxInt || opts.verbose {
		_, _ = fmt.Fprintf(opts.writer, "Processing with maxDepth: %d\n", opts.depth)
	}

	opts.processGoMod()
	opts.executeGoModGraph()
	seedNode := opts.processFile()
	printNodeWithIndentation(opts, opts.depth, 1, seedNode, "", "", 1, 1)
}

func (opts *DepTreeOpts) processFile() *Node {
	if opts.verbose {
		_, _ = fmt.Fprintf(opts.writer, "Reading go mod graph file: "+opts.goModGraphFile)
	}
	file, err := os.Open(opts.goModGraphFile)
	if err != nil {
		return nil
	}
	defer func() { _ = file.Close() }()

	// Create the seed node
	seedNode := &Node{includeVersion: opts.includeVersion}

	var nodes = make(map[string]*Node)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) != 2 {
			log.Fatalf("invalid line format: %s", line)
		}

		parent := fields[0]
		child := fields[1]

		if seedNode.Value == "" {
			seedNode.Value = parent
			nodes[parent] = seedNode
		}

		// Create nodes if they don't exist
		if _, ok := nodes[parent]; !ok {
			nodes[parent] = &Node{Value: parent, includeVersion: opts.includeVersion}
		}
		if _, ok := nodes[child]; !ok {
			nodes[child] = &Node{Value: child, Parent: nodes[parent].Val(), includeVersion: opts.includeVersion}
		}

		// Link child node to the parent only if it's not already linked
		pn := nodes[parent]
		cn := nodes[child]
		if !isChildLinked(pn, cn) {
			if pn == seedNode {
				in := opts.goModDeps[cn.Value]
				if !in {
					nodes[parent].Children = append(pn.Children, cn)
				}
			} else {
				nodes[parent].Children = append(pn.Children, cn)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf(err.Error())
	}

	return seedNode
}

func isChildLinked(parent *Node, child *Node) bool {
	for _, existingChild := range parent.Children {
		if existingChild == child {
			return true
		}
	}
	return false
}

func printNodeWithIndentation(opts *DepTreeOpts, maxDepth, depth int, node *Node, nodeIndent, childIndent string, position int, totalNodes int) {
	nonGoDeps := make([]*Node, 0)
	goDepCount := 0
	for _, child := range node.Children {
		if strings.HasPrefix(child.Val(), "golang.org") {
			goDepCount++
		} else {
			nonGoDeps = append(nonGoDeps, child)
		}
	}

	renderedNode := opts.rendered[node.Val()]
	alreadyRendered := renderedNode != ""
	openingChar := ""
	closingChar := ""
	previouslySeen := ""
	childrenMsg := ""

	if alreadyRendered {
		openingChar = " <"
		previouslySeen = "previously seen"
		closingChar = ">"
	}

	opts.rendered[node.Val()] = node.Val()
	childLen := len(nonGoDeps)

	if childLen > 0 && alreadyRendered {
		if childLen > 1 {
			childrenMsg = fmt.Sprintf(" - skipping %d children", childLen)
		} else {
			childrenMsg = " - skipping 1 child"
		}
	}

	if opts.hideSkipReason {
		_, _ = fmt.Fprintf(opts.writer, "%s%s%s\n", childIndent, nodeIndent, node.Val())
	} else {
		_, _ = fmt.Fprintf(opts.writer, "%s%s%s%s%s%s%s\n", childIndent, nodeIndent, node.Val(), openingChar, previouslySeen, childrenMsg, closingChar)
	}

	if position == totalNodes {
		childIndent += "    "
	} else {
		childIndent += " │  "
	}

	if maxDepth >= depth && renderedNode == "" {
		sort.Slice(nonGoDeps, func(i, j int) bool {
			return caseInsensitiveCompare(nonGoDeps[i].Val(), nonGoDeps[j].Val())
		})

		hasGolangDep := goDepCount > 0
		for i, child := range nonGoDeps {
			finalNode := i >= childLen-1
			if finalNode {
				if !hasGolangDep {
					nodeIndent = " └─ "
				} else {
					nodeIndent = " ├─ "
				}
			} else {
				nodeIndent = " ├─ "
			}

			if finalNode {
				if hasGolangDep {
					nodeIndent = " └─ "
					_, _ = fmt.Fprintf(opts.writer, "%s%s<skipped all [%d] golang.org* dependencies>\n", childIndent, nodeIndent, goDepCount)
					continue
				}
			}
			printNodeWithIndentation(opts, maxDepth, depth+1, child, nodeIndent, childIndent, i+1, childLen)
		}
	}
}
func (opts *DepTreeOpts) executeGoModGraph() {
	// Command to run: go mod graph
	cmd := exec.Command("go", "mod", "graph")

	// Create a file for writing the output
	outputFile, err := os.Create("./go-mod-graph.txt")
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer func() { _ = outputFile.Close() }()

	// Set the output of the command to the file
	cmd.Stdout = outputFile

	// Run the command
	err = cmd.Run()
	if err != nil {
		log.Fatalf(err.Error())
	}

	if opts.verbose {
		_, _ = fmt.Fprintf(opts.writer, "Output of go mod graph written to: ./go-mod-graph.txt\n\n")
	}
}

func caseInsensitiveCompare(a, b string) bool {
	aLower := strings.ToLower(a)
	bLower := strings.ToLower(b)

	return aLower < bLower
}

func (opts *DepTreeOpts) processGoMod() {
	if opts.verbose {
		_, _ = fmt.Fprintf(opts.writer, "Opening go.mod: "+opts.goModFile)
	}
	f, err := os.Open(opts.goModFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() { _ = f.Close() }()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		parts := strings.Split(line, " ")
		if len(parts) < 2 || strings.HasPrefix(line, "//") || line == "" {
			continue
		}

		// Check if the line is an indirect dependency
		if strings.Contains(line, "// indirect") {
			opts.goModDeps[parts[0]+"@"+parts[1]] = true
		} else {
			opts.goModDeps[parts[0]+"@"+parts[1]] = false
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	return
}
