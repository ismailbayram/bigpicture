package browser

import (
	"github.com/ismailbayram/bigpicture/internal/graph"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func ChangeDirToPythonProjectRoot() {
	err := os.Chdir("internal/browser/pyproject")
	if err != nil {
		panic(err)
	}
}

func TestPythonBrowser_GetModuleName(t *testing.T) {
	ChangeDirToPythonProjectRoot()

	browser := &PythonBrowser{
		ignoredPaths: []string{},
		tree:         nil,
	}

	moduleName := browser.getModuleName()
	assert.Equal(t, "pyproject", moduleName)
}

func TestPythonBrowser_clearNonProjectImports(t *testing.T) {

	browser := &PythonBrowser{
		ignoredPaths: []string{},
		tree:         graph.NewTree("root"),
	}
	browser.tree.Nodes["cars"] = graph.NewNode("cars", "cars", graph.Dir, []string{})
	browser.tree.Nodes["baskets"] = graph.NewNode("baskets", "./baskets", graph.Dir, []string{
		"decimal.py",
		"django/utils/translation.py",
		"cars",
	})

	browser.clearNonProjectImports()

	assert.Equal(t, 1, len(browser.tree.Nodes["baskets"].ImportRaw))
	assert.Equal(t, "cars", browser.tree.Nodes["baskets"].ImportRaw[0])
}

func TestPythonBrowser_IsIgnored(t *testing.T) {
	browser := &PythonBrowser{
		ignoredPaths: []string{"base"},
		tree:         nil,
	}

	assert.True(t, browser.isIgnored("./base/models.py"))
	assert.False(t, browser.isIgnored("./users/utils.py"))
}

func TestPythonBrowser_ParseFile(t *testing.T) {
	browser := &PythonBrowser{
		ignoredPaths: []string{},
		tree:         nil,
		moduleName:   "pyproject",
	}
	parentNode := graph.NewNode("pyproject", "./", graph.Dir, nil)
	node := browser.parseFile("baskets/service.py", parentNode)
	assert.NotNil(t, node)
	assert.Equal(t, "baskets/service.py", node.Path)
	assert.Equal(t, "service.py", node.PackageName)

	assert.Equal(t, "/decimal.py", node.ImportRaw[0])
	assert.Equal(t, "/django/utils/translation.py", node.ImportRaw[1])
	assert.Equal(t, "/django/db/models.py", node.ImportRaw[2])
	assert.Equal(t, "/django/db/transaction.py", node.ImportRaw[3])
	assert.Equal(t, "/baskets/models.py", node.ImportRaw[4])
	assert.Equal(t, "/baskets/enums.py", node.ImportRaw[5])
	assert.Equal(t, "/baskets/exceptions.py", node.ImportRaw[6])
	assert.Equal(t, "/cars/exceptions.py", node.ImportRaw[7])
	assert.Equal(t, "/cars", node.ImportRaw[8])
}

func TestPythonBrowser_Browse(t *testing.T) {
	browser := NewBrowser(LangPy, graph.NewTree("root"), []string{}).(*PythonBrowser)

	browser.Browse(".")
	assert.Equal(t, "pyproject", browser.moduleName)
	assert.NotEqual(t, 1, len(browser.tree.Nodes))
}

func TestPythonBrowser_browse(t *testing.T) {
	browser := NewBrowser(LangPy, graph.NewTree("root"), []string{}).(*PythonBrowser)

	browser.browse("base/", browser.tree.Root)

	assert.Equal(t, 6, len(browser.tree.Nodes))
}
