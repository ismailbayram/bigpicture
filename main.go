package main

import (
	"bytes"
	"fmt"
	"github.com/ismailbayram/bigpicture/internal/browser"
	"github.com/ismailbayram/bigpicture/internal/graph"
	"os"
	"strings"
)

func main() {
	file, err := os.ReadFile("go.mod")
	if err != nil {
		panic(err)
	}
	firstLine := string(file[:bytes.IndexByte(file, '\n')])
	moduleName := strings.Split(firstLine, " ")[1]
	rootNode := graph.NewNode(moduleName, ".", nil, graph.Dir)

	err = browser.Browse(".", moduleName, rootNode)
	if err != nil {
		panic(err)
	}

	PrintTree(rootNode, 0)
}

func PrintTree(node *graph.Node, depth int) {
	for _, child := range node.Children {
		fmt.Println(strings.Repeat(" ", depth), child.PackageName, child.Path, child.Type)
		PrintTree(child, depth+1)
	}
}
