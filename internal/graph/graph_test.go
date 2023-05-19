package graph

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTree(t *testing.T) {
	root := NewNode("package", "path/dir", nil, Dir, nil)
	tree := NewTree(root)
	assert.Equal(t, root, tree.Root)
	assert.Equal(t, 1, len(tree.Nodes))
	assert.Equal(t, root, tree.Nodes["path/dir"])
}

func TestTree_ConvertImportRaw(t *testing.T) {
	root := NewNode("main", "path", nil, Dir, []string{"path/node1", "path/node2"})
	tree := NewTree(root)

	node1 := NewNode("package1", "path/node1", root, Dir, []string{})
	tree.Nodes["path/node1"] = node1
	node1_file := NewNode("package1", "path/node1/file", node1, File, []string{"path/node2"})
	tree.Nodes["path/node1/file"] = node1_file
	node2 := NewNode("package2", "path/node2", root, Dir, []string{"path/node3"})
	tree.Nodes["path/node2"] = node2
	node3 := NewNode("package3", "path/node3", root, Dir, []string{})
	tree.Nodes["path/node3"] = node3

	tree.ConvertImportRaw()
	assert.Equal(t, 2, len(root.Imports))
	assert.Contains(t, root.Imports, node1)
	assert.Contains(t, root.Imports, node2)

	assert.Equal(t, 1, len(node1.Imports))
	assert.Contains(t, node1.Imports, node2)
	assert.Equal(t, 1, len(node1_file.Imports))

	assert.Equal(t, 1, len(node2.Imports))
	assert.Contains(t, node2.Imports, node3)

	assert.Equal(t, 0, len(node3.Imports))
}

func TestNewNode(t *testing.T) {
	node := NewNode("package", "path/dir", nil, Dir, nil)
	assert.Equal(t, "package", node.PackageName)
	assert.Equal(t, "path/dir", node.Path)
	assert.Equal(t, 0, len(node.Imports))
	assert.Nilf(t, node.Parent, "Parent should be nil.")
	assert.Equal(t, Dir, node.Type)

	node = NewNode("package", "./path/dir", nil, Dir, []string{"import1", "import2"})
	assert.Equal(t, "package", node.PackageName)
	assert.Equal(t, "/path/dir", node.Path)
	assert.Equal(t, 2, len(node.importRaw))
}

func TestNode_AddChild(t *testing.T) {
	parentNode := NewNode("package", "parent", nil, Dir, nil)
	childNode := NewNode("child package", "parent/child", parentNode, Dir, nil)
	parentNode.AddChild(childNode)
	assert.Equal(t, 1, len(parentNode.Children))
	assert.Equal(t, childNode, parentNode.Children[0])
}
