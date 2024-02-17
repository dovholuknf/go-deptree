package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Node struct {
	Value    string
	Children []*Node
}

func main() {
	fmt.Println("Version 1")
	// Define a maxDepth flag with a default value of 0
	maxDepth := flag.Int("maxDepth", 1, "Maximum depth for processing")

	// Parse the command-line arguments
	flag.Parse()

	// Access the value of the maxDepth flag
	depth := *maxDepth

	// Validate the value (optional)
	if depth < 0 {
		fmt.Println("Error: maxDepth must be a non-negative integer.")
		return
	}

	// Your program logic goes here, using the 'depth' variable

	// Example: Print the maxDepth value
	fmt.Printf("Processing with maxDepth: %d\n", depth)
	executeGoModGraph()
	filePath := "./go-mod-graph.txt"
	seedNode, err := processFile(filePath, getCurrentModuleName())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	printNodeWithIndentation(*maxDepth, 1, seedNode, "", "", 1, 1)
}

func processFile(filePath string, seedValue string) (*Node, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	nodes := make(map[string]*Node)

	// Create the seed node
	seedNode := &Node{Value: seedValue}
	nodes[seedValue] = seedNode

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) != 2 {
			return nil, fmt.Errorf("invalid line format: %s", line)
		}

		parent := fields[0]
		child := fields[1]

		// Create nodes if they don't exist
		if _, ok := nodes[parent]; !ok {
			nodes[parent] = &Node{Value: parent}
		}
		if _, ok := nodes[child]; !ok {
			nodes[child] = &Node{Value: child}
		}

		// Link child node to the parent only if it's not already linked
		if !isChildLinked(nodes[parent], nodes[child]) {
			nodes[parent].Children = append(nodes[parent].Children, nodes[child])
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
	childLen := len(node.Children)

	fmt.Print("%s%s%s", childIndent, nodeIndent, node.Value)
	if strings.HasPrefix(node.Value, "golang.org") {
		fmt.Printf(" <skipping %d children from golang.org>\n", childLen)
	}

	if position == totalNodes {
		childIndent += "    "
	} else {
		childIndent += "│   "
	}
	if maxDepth > depth {
		for i, child := range node.Children {
			if i == childLen-1 {
				nodeIndent = "└── "
			} else {
				nodeIndent = "├── "
			}
			printNodeWithIndentation(maxDepth, depth+1, child, nodeIndent, childIndent, i+1, childLen)
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
	defer outputFile.Close()

	// Set the output of the command to the file
	cmd.Stdout = outputFile

	// Run the command
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error running command:", err)
		return
	}

	fmt.Println("Command successfully executed. Output written to ./go-mod-graph.txt")
}

func getCurrentModuleName() string {
	cmd := exec.Command("go", "list", "-m")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	// Convert the output to a string and trim any whitespace
	moduleName := strings.TrimSpace(string(output))
	return moduleName
}
