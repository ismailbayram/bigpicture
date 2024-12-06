package browser

import (
	"github.com/ismailbayram/bigpicture/internal/graph"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func ChangeDirToPythonProjectRoot() {
	err := os.Chdir("../pyproject")
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

func TestPythonBrowser_ParseFile(t *testing.T) {
	browser := &PythonBrowser{
		ignoredPaths: []string{},
		tree:         nil,
		moduleName:   "pyproject",
		rootDir:      "/",
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

	assert.Equal(t, 216, node.LineCount)

	assert.Equal(t, 11, len(node.Functions))
	funcs := map[string]int{
		"get_or_create_basket":    16,
		"apply_discounts":         15,
		"_check_basket_items":     17,
		"add_basket_item":         22,
		"clean_discounts":         7,
		"clean_basket":            7,
		"delete_basket_item":      7,
		"complete_basket":         19,
		"_create_quotation_items": 14,
		"check_discount_code":     24,
		"apply_campaigns":         29,
	}
	for _, f := range node.Functions {
		assert.Equal(t, funcs[f.Name], f.LineCount, f.Name)
	}

}

func TestPythonBrowser_Browse(t *testing.T) {
	browser := NewBrowser(LangPy, graph.NewTree("root"), []string{}, "/").(*PythonBrowser)

	browser.Browse(".")
	assert.Equal(t, "pyproject", browser.moduleName)
	assert.NotEqual(t, 1, len(browser.tree.Nodes))
}

func TestPythonBrowser_browse(t *testing.T) {
	browser := NewBrowser(LangPy, graph.NewTree("root"), []string{}, "/").(*PythonBrowser)

	browser.browse("base/", browser.tree.Root)

	assert.Equal(t, 6, len(browser.tree.Nodes))
	assert.Equal(t, browser.tree.Root.LineCount, 191)
}
