package browser

import (
	"fmt"
	"github.com/ismailbayram/bigpicture/internal/graph"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func ChangeDirToJavaProjectRoot() {
	err := os.Chdir("internal/browser/javaproject")
	if err != nil {
		panic(err)
	}
}

func TestJavaBrowser_GetModuleName(t *testing.T) {
	ChangeDirToJavaProjectRoot()

	browser := &JavaBrowser{
		ignoredPaths: []string{},
		tree:         nil,
	}

	moduleName := browser.getModuleName()
	assert.Equal(t, "javaproject", moduleName)
}

func TestJavaBrowser_clearNonProjectImports(t *testing.T) {

	browser := &JavaBrowser{
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

func TestJavaBrowser_ParseFile(t *testing.T) {
	browser := &JavaBrowser{
		ignoredPaths: []string{},
		tree:         nil,
		moduleName:   "javaproject",
	}
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(dir)
	parentNode := graph.NewNode("javaproject", "./", graph.Dir, nil)
	node := browser.parseFile("src/com/shashi/service/impl/TrainServiceImpl.java", parentNode)
	assert.NotNil(t, node)

	assert.Equal(t, "src/com/shashi/service/impl/TrainServiceImpl.java", node.Path)
	assert.Equal(t, "TrainServiceImpl.java", node.PackageName)
	assert.Equal(t, "/com/shashi/beans/TrainBean.java", node.ImportRaw[0])
	assert.Equal(t, "/com/shashi/beans/TrainException.java", node.ImportRaw[1])
	assert.Equal(t, "/com/shashi/constant/ResponseCode.java", node.ImportRaw[2])
	assert.Equal(t, "/com/shashi/service/TrainService.java", node.ImportRaw[3])
	assert.Equal(t, "/com/shashi/utility/DBUtil.java", node.ImportRaw[4])

	assert.Equal(t, 170, node.LineCount)

	assert.Equal(t, 6, len(node.Functions))
	funcs := map[string]int{
		"addTrain":                 20,
		"deleteTrainById":          15,
		"updateTrain":              20,
		"getTrainById":             22,
		"getAllTrains":             24,
		"getTrainsBetweenStations": 27,
	}
	for _, f := range node.Functions {
		assert.Equal(t, funcs[f.Name], f.LineCount, f.Name)
	}
}

func TestJavaBrowser_Browse(t *testing.T) {
	browser := NewBrowser(LangJava, graph.NewTree("root"), []string{}).(*JavaBrowser)

	browser.Browse(".")
	assert.Equal(t, "javaproject", browser.moduleName)
	assert.NotEqual(t, 1, len(browser.tree.Nodes))
}

func TestJavaBrowser_browse(t *testing.T) {
	browser := NewBrowser(LangJava, graph.NewTree("root"), []string{}).(*JavaBrowser)

	browser.browse("src/com/shashi/service", browser.tree.Root)

	assert.Equal(t, 8, len(browser.tree.Nodes))
}
