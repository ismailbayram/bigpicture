package browser

import (
	"fmt"
	"github.com/ismailbayram/bigpicture/internal/graph"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func Browse(path string, moduleName string, parentNode *graph.Node) error {
	entries, err := os.ReadDir(path)

	if err != nil {
		return err
	}

	for _, e := range entries {
		fName := e.Name()
		if e.IsDir() && !strings.Contains(fName, ".") {
			dirPath := fmt.Sprintf("%s/%s", path, fName)
			node := graph.NewNode(fName, dirPath, parentNode, graph.Dir)
			parentNode.AddChild(node)
			if err := Browse(dirPath, moduleName, node); err != nil {
				return err
			}
		} else if strings.HasSuffix(fName, ".go") {
			return parseFile(fName, path, moduleName, parentNode)
		}
	}
	return nil
}

func parseFile(fName string, path, moduleName string, parentNode *graph.Node) error {
	filePath := fmt.Sprintf("%s/%s", path, fName)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ImportsOnly)
	if err != nil {
		return err
	}

	for _, s := range f.Imports {
		if strings.Contains(s.Path.Value, moduleName) {
			fmt.Println(s.Path.Value)
		}
	}
	fmt.Println("----")

	node := graph.NewNode(fName, filePath, parentNode, graph.File)
	parentNode.AddChild(node)
	parentNode.PackageName = f.Name.Name
	node.PackageName = f.Name.Name

	return nil
}
