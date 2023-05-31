package browser

import (
	"fmt"
	"github.com/ismailbayram/bigpicture/internal/graph"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

type Browser struct {
	ignoredPaths []string
	moduleName   string
	tree         *graph.Tree
}

func NewBrowser(moduleName string, tree *graph.Tree, ignoredPaths []string) *Browser {
	return &Browser{
		ignoredPaths: ignoredPaths,
		moduleName:   moduleName,
		tree:         tree,
	}
}

func (b *Browser) isIgnored(entryPath string) bool {
	entryPath = entryPath[2:]
	for _, path := range b.ignoredPaths {
		if strings.HasPrefix(entryPath, path) {
			return true
		}
	}
	return false
}

func (b *Browser) Browse(parentPath string) {
	b.browse(parentPath, b.tree.Root)
}

func (b *Browser) browse(parentPath string, parentNode *graph.Node) {
	entries, err := os.ReadDir(parentPath)

	if err != nil {
		panic(err)
	}

	for _, e := range entries {
		fName := e.Name()
		path := fmt.Sprintf("%s/%s", parentPath, fName)
		if b.isIgnored(path) {
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

func (b *Browser) parseFile(path string, parentNode *graph.Node) *graph.Node {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
	if err != nil {
		panic(err)
	}

	var imports []string
	for _, s := range f.Imports {
		if strings.Contains(s.Path.Value, b.moduleName) {
			_path := strings.Split(strings.Trim(s.Path.Value, "\""), b.moduleName)[1]
			imports = append(imports, _path)
		}
	}

	node := graph.NewNode(f.Name.Name, path, graph.File, imports)
	parentNode.PackageName = f.Name.Name
	node.PackageName = f.Name.Name

	return node
}

//func TestConfiguration_IsIgnored(t *testing.T) {
//	defer os.Remove(FileName)
//
//	cfg := Init()
//	cfg.IgnoredPaths = []string{"vendor", "web"}
//
//	assert.True(t, cfg.IsIgnored("./web/something"))
//	assert.False(t, cfg.IsIgnored("./cmd"))
//}
