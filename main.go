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
	moduleName := getModuleName()
	rootNode := graph.NewNode(moduleName, ".", nil, graph.Dir, nil)
	tree := graph.NewTree(rootNode)

	if err := browser.Browse(".", moduleName, rootNode, tree); err != nil {
		panic(err)
	}

	tree.ConvertImportRaw()
	PrintTree(tree.Root, 0)
}

func PrintTree(node *graph.Node, depth int) {

	for _, child := range node.Children {
		fmt.Println(strings.Repeat(" ", depth), child.PackageName, child.Path, child.Type)
		fmt.Println(strings.Repeat(" ", depth), child.Imports)
		PrintTree(child, depth+1)
	}
}

func getModuleName() string {
	file, err := os.ReadFile("go.mod")
	if err != nil {
		panic(err)
	}
	firstLine := string(file[:bytes.IndexByte(file, '\n')])
	return strings.Split(firstLine, " ")[1]
}
