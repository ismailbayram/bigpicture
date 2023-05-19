package browser

import (
	"fmt"
	"github.com/ismailbayram/bigpicture/internal/graph"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func Browse(parentPath string, moduleName string, parentNode *graph.Node, tree *graph.Tree) error {
	entries, err := os.ReadDir(parentPath)

	if err != nil {
		return err
	}

	for _, e := range entries {
		fName := e.Name()
		path := fmt.Sprintf("%s/%s", parentPath, fName)

		if e.IsDir() && !strings.Contains(fName, ".") {
			node := graph.NewNode(fName, path, parentNode, graph.Dir, nil)
			parentNode.AddChild(node)
			if err := Browse(path, moduleName, node, tree); err != nil {
				return err
			}
		} else if strings.HasSuffix(fName, ".go") {
			if err := parseFile(path, moduleName, parentNode); err != nil {
				return err
			}
		}
	}

	for _, node := range parentNode.Children {
		tree.Nodes[node.Path] = node
	}

	return nil
}

func parseFile(path string, moduleName string, parentNode *graph.Node) error {

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
	if err != nil {
		return err
	}

	var imports []string
	for _, s := range f.Imports {
		if strings.Contains(s.Path.Value, moduleName) {
			_path := strings.Split(strings.Trim(s.Path.Value, "\""), moduleName)[1]
			imports = append(imports, _path)
		}
	}

	node := graph.NewNode(f.Name.Name, path, parentNode, graph.File, imports)
	parentNode.AddChild(node)
	parentNode.PackageName = f.Name.Name
	node.PackageName = f.Name.Name

	return nil
}
