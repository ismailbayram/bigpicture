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

	node1 := NewNode("package1", "/node1", Dir, []string{})
	tree.Nodes["/node1"] = node1
	node1File := NewNode("package1", "/node1/file", File, []string{"/node2"})
	tree.Nodes["/node1/file"] = node1File
	node2 := NewNode("package2", "/node2", Dir, []string{"/node3"})
	tree.Nodes["/node2"] = node2
	node3 := NewNode("package3", "/node3", Dir, []string{})
	tree.Nodes["/node3"] = node3

	tree.GenerateLinks()

	assert.Equal(t, 3, len(tree.Links))

	assert.Equal(t, node1File, tree.Links[0].From)
	assert.Equal(t, node2, tree.Links[0].To)
	assert.True(t, tree.Links[0].IsVisible)

	assert.Equal(t, node1, tree.Links[1].From)
	assert.Equal(t, node2, tree.Links[1].To)
	assert.False(t, tree.Links[1].IsVisible)

	assert.Equal(t, node2, tree.Links[2].From)
	assert.Equal(t, node3, tree.Links[2].To)
	assert.True(t, tree.Links[2].IsVisible)
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
	assert.Equal(t, node.Parent, "/path")
	assert.Equal(t, 2, len(node.ImportRaw))

	node = NewNode("package", "/path", Dir, []string{})
	assert.Equal(t, ".", node.Parent)
}

func TestNode_FileName(t *testing.T) {
	node := NewNode("package", "path/dir", Dir, nil)
	assert.Equal(t, "dir", node.FileName())

	node = NewNode("package", "path/dir/file.go", File, nil)
	assert.Equal(t, "file", node.FileName())
}

func TestNode_ToJSON(t *testing.T) {
	node := NewNode("package", "path/dir", Dir, nil)
	node.Functions = []Function{{Name: "test", LineCount: 10}}
	assert.Equal(
		t,
		`{"package_name":"package","path":"path/dir","parent":"path","type":1,"line_count":0,"functions":[{"name":"test","line_count":10}],"instability":null}`,
		node.ToJSON(),
	)
}

func TestNewLink(t *testing.T) {
	node1 := NewNode("package1", "path/node1", Dir, []string{})
	node2 := NewNode("package2", "path/node2", Dir, []string{})
	link := NewLink(node1, node2, true)
	assert.Equal(t, node1, link.From)
	assert.Equal(t, node2, link.To)
	assert.True(t, link.IsVisible)
}

func TestLink_ToJSON(t *testing.T) {
	node1 := NewNode("package1", "path/node1", Dir, []string{})
	node2 := NewNode("package2", "path/node2", Dir, []string{})
	link := NewLink(node1, node2, true)
	assert.Equal(t, "{\"from\":\"path/node1\",\"to\":\"path/node2\",\"is_visible\":true}", link.ToJSON())
}
