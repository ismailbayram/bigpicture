package graph

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTree(t *testing.T) {
	root := NewNode("package", "path/dir", Dir, nil)
	tree := NewTree(root)
	assert.Equal(t, root, tree.Root)
	assert.Equal(t, 1, len(tree.Nodes))
	assert.Equal(t, root, tree.Nodes["path/dir"])
}

func TestTree_GenerateLinks(t *testing.T) {
	root := NewNode("main", "path", Dir, []string{"path/node1", "path/node2"})
	tree := NewTree(root)

	node1 := NewNode("package1", "path/node1", Dir, []string{})
	tree.Nodes["path/node1"] = node1
	node1_file := NewNode("package1", "path/node1/file", File, []string{"path/node2"})
	tree.Nodes["path/node1/file"] = node1_file
	node2 := NewNode("package2", "path/node2", Dir, []string{"path/node3"})
	tree.Nodes["path/node2"] = node2
	node3 := NewNode("package3", "path/node3", Dir, []string{})
	tree.Nodes["path/node3"] = node3

	tree.GenerateLinks()

	assert.Equal(t, 5, len(tree.Links))

	assert.Equal(t, root, tree.Links[0].From)
	assert.Equal(t, node1, tree.Links[0].To)

	assert.Equal(t, root, tree.Links[1].From)
	assert.Equal(t, node2, tree.Links[1].To)

	assert.Equal(t, node1_file, tree.Links[2].From)
	assert.Equal(t, node2, tree.Links[2].To)

	assert.Equal(t, node1, tree.Links[3].From)
	assert.Equal(t, node2, tree.Links[3].To)

	assert.Equal(t, node2, tree.Links[4].From)
	assert.Equal(t, node3, tree.Links[4].To)
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
	assert.Equal(t, 2, len(node.importRaw))
}
