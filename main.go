package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"math"
	"os"
	"os/exec"
	"sort"
	"strings"
	"testing"
	"time"
)

const ver = "1.0.9"

func TestSomething(t *testing.T) {

	// assert equality
	assert.Equal(t, 123, 123, "they should be equal")
}

var verbose = false
var includeVersion = false

// var hideSkipReason = false
var rendered = make(map[string]string)
var directDeps = make(map[string]bool)

//var indirectDeps = make(map[string]bool)

type Node struct {
	Value    string
	Children []*Node
	Parent   string
	Indirect bool
}

func (n *Node) Val() string {
	if includeVersion {
		return n.Value
	} else {
		return strings.Split(n.Value, "@")[0]
	}
}

func main() {
	time.Sleep(time.Second)
	maxDepth := flag.Int("maxDepth", math.MaxInt, "Maximum depth for processing")
	var versionFlag = false
	flag.BoolVar(&versionFlag, "version", false, "Print the version")
	flag.BoolVar(&verbose, "verbose", false, "Print additional output")
	//	flag.BoolVar(&hideSkipReason, "hideSkipReason", false, "Suppresses the reason for skipping child dependencies")

	flag.BoolVar(&includeVersion, "includeVersion", false, "Prints the version of the dependency too")
	flag.Parse()

	if versionFlag {
		fmt.Println(ver)
		return
	}

	depth := *maxDepth

	if depth < 1 {
		fmt.Println("maxDepth cannot be < 1, using 1 for maxDepth")
		depth = 1
	}
	if depth < math.MaxInt || verbose {
		fmt.Printf("Processing with maxDepth: %d\n", depth)
	}

	goModFile := "./go.mod"
	//goModFile = "/home/cd/git/github/dovholuknf/edgex/go-mod-bootstrap/go.mod"
	processGoMod(goModFile)

	executeGoModGraph()
	filePath := "./go-mod-graph.txt"
	//filePath = "/home/cd/git/github/dovholuknf/edgex/go-mod-bootstrap/go-mod-graph.txt"
	seedNode, err := processFile(filePath) // getCurrentModuleName())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	printNodeWithIndentation(depth, 1, seedNode, "", "", 1, 1)
}

func processFile(goModGraphFile string) (*Node, error) {
	if verbose {
		fmt.Println("Reading go mod graph file: " + goModGraphFile)
	}
	file, err := os.Open(goModGraphFile)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	// Create the seed node
	seedNode := &Node{}

	var nodes = make(map[string]*Node)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) != 2 {
			return nil, fmt.Errorf("invalid line format: %s", line)
		}

		parent := fields[0]
		child := fields[1]

		if seedNode.Value == "" {
			seedNode.Value = parent
			nodes[parent] = seedNode
		}

		// Create nodes if they don't exist
		if _, ok := nodes[parent]; !ok {
			nodes[parent] = &Node{Value: parent}
		}
		if _, ok := nodes[child]; !ok {
			nodes[child] = &Node{Value: child, Parent: nodes[parent].Val()}
		}

		// Link child node to the parent only if it's not already linked
		pn := nodes[parent]
		cn := nodes[child]
		if !isChildLinked(pn, cn) {
			if pn == seedNode {
				in := directDeps[cn.Value]
				if in {
					//nodes[parent].Children = append(pn.Children, cn)
				} else {
					nodes[parent].Children = append(pn.Children, cn)
				}
			} else {
				nodes[parent].Children = append(pn.Children, cn)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return seedNode, nil
}

func isChildLinked(parent *Node, child *Node) bool {
	for _, existingChild := range parent.Children {
		if existingChild == child {
			return true
		}
	}
	return false
}

func printNodeWithIndentation(maxDepth, depth int, node *Node, nodeIndent, childIndent string, position int, totalNodes int) {
	renderedNode := rendered[node.Val()]
	alreadyRendered := renderedNode != ""
	//msgFmt := "%s%s%s"
	openingChar := ""
	closingChar := ""
	previouslySeen := ""
	chillen := ""

	if alreadyRendered {
		openingChar = " <"
		previouslySeen = "previously seen"
		closingChar = ">"
	}

	rendered[node.Val()] = node.Val()
	childLen := len(node.Children)

	if childLen > 0 && alreadyRendered {
		if childLen > 1 {
			chillen = fmt.Sprintf(" - skipping %d children", childLen)
		} else {
			chillen = " - skipping 1 child"
		}
	}

	fmt.Printf("%s%s%s%s%s%s%s\n", childIndent, nodeIndent, node.Val(), openingChar, previouslySeen, chillen, closingChar)

	if position == totalNodes {
		childIndent += "    "
	} else {
		childIndent += "│   "
	}

	if maxDepth >= depth && renderedNode == "" {
		sort.Slice(node.Children, func(i, j int) bool {
			return caseInsensitiveCompare(node.Children[i].Val(), node.Children[j].Val())
		})

		hasGolangDep := false
		for i, child := range node.Children {
			if strings.HasPrefix(child.Val(), "golang.org") {
				hasGolangDep = true //make a note of golang dep and carry on
				continue
			}
			if child.Val() == "github.com/tklauser/numcpus" {
				println("XXXX")
				println("XXXX")
				println("XXXX")
				println("XXXX")
			}

			finalNode := i >= childLen-1
			if finalNode {
				if !hasGolangDep {
					nodeIndent = "└── "
				} else {
					nodeIndent = "├── "
				}
			} else {
				nodeIndent = "├── "
			}

			if finalNode {
				fmt.Printf("%s%s%s%s%s%s%s\n", childIndent, "└── ", "<skipped all golang.org* dependencies>", "", "", "", "")
			} else {
				printNodeWithIndentation(maxDepth, depth+1, child, nodeIndent, childIndent, i+1, childLen)
			}
		}
	}
}
func executeGoModGraph() {
	// Command to run: go mod graph
	cmd := exec.Command("go", "mod", "graph")

	// Create a file for writing the output
	outputFile, err := os.Create("./go-mod-graph.txt")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer func() { _ = outputFile.Close() }()

	// Set the output of the command to the file
	cmd.Stdout = outputFile

	// Run the command
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error running command:", err)
		return
	}

	if verbose {
		fmt.Printf("Output of go mod graph written to: ./go-mod-graph.txt\n\n")
	}
}

func caseInsensitiveCompare(a, b string) bool {
	aLower := strings.ToLower(a)
	bLower := strings.ToLower(b)

	return aLower < bLower
}

func processGoMod(file string) {
	if verbose {
		fmt.Println("Opening go.mod: " + file)
	}
	f, err := os.Open(file)
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
			directDeps[parts[0]+"@"+parts[1]] = true
		} else {
			directDeps[parts[0]+"@"+parts[1]] = false
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	return
}
