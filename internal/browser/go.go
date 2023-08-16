package browser

import (
	"bytes"
	"fmt"
	"github.com/ismailbayram/bigpicture/internal/graph"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

type GoBrowser struct {
	ignoredPaths []string
	moduleName   string
	tree         *graph.Tree
}

func (b *GoBrowser) Browse(parentPath string) {
	b.moduleName = b.getModuleName()
	b.browse(parentPath, b.tree.Root)
}

func (b *GoBrowser) getModuleName() string {
	file, err := os.ReadFile("go.mod")
	if err != nil {
		panic(err)
	}
	firstLine := string(file[:bytes.IndexByte(file, '\n')])
	return strings.Split(firstLine, " ")[1]
}

func (b *GoBrowser) browse(parentPath string, parentNode *graph.Node) {
	entries, err := os.ReadDir(parentPath)

	if err != nil {
		panic(err)
	}

	for _, e := range entries {
		fName := e.Name()
		path := fmt.Sprintf("%s/%s", parentPath, fName)
		if isIgnored(b.ignoredPaths, path) {
			continue
		}

		if e.IsDir() && !strings.Contains(fName, ".") {
			node := graph.NewNode(fName, path, graph.Dir, nil)
			b.tree.Nodes[node.Path] = node
			b.browse(path, node)
		} else if strings.HasSuffix(fName, ".go") {
			node := b.parseFile(path, parentNode)
			b.tree.Nodes[node.Path] = node
		}
	}
}

func (b *GoBrowser) parseFile(path string, parentNode *graph.Node) *graph.Node {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	fileContent := string(file)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", fileContent, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	// extract imports
	var imports []string
	for _, s := range f.Imports {
		if strings.Contains(s.Path.Value, b.moduleName) {
			_path := strings.Split(strings.Trim(s.Path.Value, "\""), b.moduleName)[1]
			imports = append(imports, _path)
		}
	}

	// extract functions
	functions := make([]graph.Function, 0)
	for _, decl := range f.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			functions = append(functions, graph.Function{
				Name:      fn.Name.Name,
				LineCount: fset.Position(fn.End()).Line - fset.Position(fn.Pos()).Line - 2,
			})
		}
	}

	node := graph.NewNode(f.Name.Name, path, graph.File, imports)
	node.Functions = functions
	parentNode.PackageName = f.Name.Name
	node.PackageName = f.Name.Name
	node.LineCount = strings.Count(fileContent, "\n")

	return node
}
