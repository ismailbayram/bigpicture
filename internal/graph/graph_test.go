package graph

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTree(t *testing.T) {
	tree := NewTree("moduleName")
	assert.Equal(t, "moduleName", tree.Root.PackageName)
	assert.Equal(t, 1, len(tree.Nodes))
	assert.Equal(t, tree.Root, tree.Nodes["."])
}

func TestTree_GenerateLinks(t *testing.T) {
	tree := NewTree("moduleName")
	tree.Root.ImportRaw = []string{"path/node1", "path/node2"}

	node1 := NewNode("package1", "path/node1", Dir, []string{})
	tree.Nodes["path/node1"] = node1
	node1_file := NewNode("package1", "path/node1/file", File, []string{"path/node2"})
	tree.Nodes["path/node1/file"] = node1_file
	node2 := NewNode("package2", "path/node2", Dir, []string{"path/node3"})
	tree.Nodes["path/node2"] = node2
	node3 := NewNode("package3", "path/node3", Dir, []string{})
	tree.Nodes["path/node3"] = node3

	tree.GenerateLinks()

	assert.Equal(t, 4, len(tree.Links))

	assert.Equal(t, tree.Root, tree.Links[0].From)
	assert.Equal(t, node1, tree.Links[0].To)

	assert.Equal(t, tree.Root, tree.Links[1].From)
	assert.Equal(t, node2, tree.Links[1].To)

	assert.Equal(t, node1_file, tree.Links[2].From)
	assert.Equal(t, node2, tree.Links[2].To)

	//assert.Equal(t, node1, tree.Links[3].From)
	//assert.Equal(t, node2, tree.Links[3].To)

	assert.Equal(t, node2, tree.Links[3].From)
	assert.Equal(t, node3, tree.Links[3].To)
}

func TestNewNode(t *testing.T) {
	node := NewNode("package", "path/dir", Dir, nil)
	assert.Equal(t, "package", node.PackageName)
	assert.Equal(t, "path/dir", node.Path)
	assert.Equal(t, node.Parent, "path")
	assert.Equal(t, Dir, node.Type)

	node = NewNode("package", "./path/dir", Dir, []string{"import1", "import2"})
	assert.Equal(t, "package", node.PackageName)
	assert.Equal(t, "/path/dir", node.Path)
	assert.Equal(t, 2, len(node.ImportRaw))
}

func TestNode_ToJSON(t *testing.T) {
	node := NewNode("package", "path/dir", Dir, nil)
	assert.Equal(t, "{\"package_name\":\"package\",\"path\":\"path/dir\",\"parent\":\"path\",\"type\":1}", node.ToJSON())
}

func TestNewLink(t *testing.T) {
	node1 := NewNode("package1", "path/node1", Dir, []string{})
	node2 := NewNode("package2", "path/node2", Dir, []string{})
	link := NewLink(node1, node2)
	assert.Equal(t, node1, link.From)
	assert.Equal(t, node2, link.To)
}

func TestLink_ToJSON(t *testing.T) {
	node1 := NewNode("package1", "path/node1", Dir, []string{})
	node2 := NewNode("package2", "path/node2", Dir, []string{})
	link := NewLink(node1, node2)
	assert.Equal(t, "{\"from\":\"path/node1\",\"to\":\"path/node2\"}", link.ToJSON())
}
