package browser

import (
	"github.com/ismailbayram/bigpicture/internal/graph"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGoBrowser_GetModuleName(t *testing.T) {
	os.Chdir("../..")

	browser := &GoBrowser{
		ignoredPaths: []string{},
		tree:         nil,
	}

	moduleName := browser.getModuleName()
	assert.Equal(t, "github.com/ismailbayram/bigpicture", moduleName)
}

func TestGoBrowser_ParseFile(t *testing.T) {
	browser := &GoBrowser{
		ignoredPaths: []string{},
		tree:         nil,
		moduleName:   "github.com/ismailbayram/bigpicture",
	}
	parentNode := graph.NewNode("bigpicture", "./", graph.Dir, nil)
	node := browser.parseFile("main.go", parentNode)
	assert.NotNil(t, node)
	assert.Equal(t, "main.go", node.Path)
	assert.Equal(t, "main", node.PackageName)
	assert.Equal(t, parentNode.PackageName, node.PackageName)

	assert.Equal(t, "/internal/browser", node.ImportRaw[0])
	assert.Equal(t, "/internal/config", node.ImportRaw[1])
	assert.Equal(t, "/internal/graph", node.ImportRaw[2])
	assert.Equal(t, "/internal/server", node.ImportRaw[3])
}

func TestGoBrowser_Browse(t *testing.T) {
	browser := NewBrowser(LangGo, graph.NewTree("root"), []string{}, "/").(*GoBrowser)

	browser.Browse(".")
	assert.Equal(t, "github.com/ismailbayram/bigpicture", browser.moduleName)
	assert.NotEqual(t, 1, len(browser.tree.Nodes))
}

func TestGoBrowser_browse(t *testing.T) {
	browser := NewBrowser(
		LangGo,
		graph.NewTree("root"),
		[]string{"internal/browser/pyproject", "internal/browser/javaproject"},
		"/",
	).(*GoBrowser)

	browser.browse("./internal/browser", browser.tree.Root)

	assert.Equal(t, 9, len(browser.tree.Nodes))
}
