package graph

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNode(t *testing.T) {
	node := NewNode("package", "path/dir", nil, Dir)
	assert.Equal(t, "package", node.PackageName)
	assert.Equal(t, "path/dir", node.Path)
	assert.Equal(t, 0, len(node.Imports))
	assert.Nilf(t, node.Parent, "Parent should be nil.")
}

func TestNode_AddImports(t *testing.T) {
	parentNode := NewNode("package", "parent", nil, Dir)
	importingNode := NewNode("importing package", "parent/importing", parentNode, Dir)
	importedNode := NewNode("other package", "path/dir", nil, Dir)

	importingNode.AddImport(importedNode)
	assert.Equal(t, 1, len(importingNode.Imports))
	assert.Equal(t, 1, len(parentNode.Imports))
	assert.Equal(t, parentNode.Imports[0], importingNode.Imports[0])
	edge := importingNode.Imports[0]
	assert.Equal(t, importedNode, edge.To)
	assert.Equal(t, importingNode, edge.From)
}

func TestNewEdge(t *testing.T) {
	nodeFrom := NewNode("package", "path/dir", nil, Dir)
	nodeTo := NewNode("package", "path/dir", nil, Dir)
	edge := NewEdge(nodeFrom, nodeTo)
	assert.Equal(t, nodeFrom, edge.From)
	assert.Equal(t, nodeTo, edge.To)
}

func TestNode_AddChild(t *testing.T) {
	parentNode := NewNode("package", "parent", nil, Dir)
	childNode := NewNode("child package", "parent/child", parentNode, Dir)
	parentNode.AddChild(childNode)
	assert.Equal(t, 1, len(parentNode.Children))
	assert.Equal(t, childNode, parentNode.Children[0])
}
