package graph

import "strings"

type Type uint

const (
	Dir Type = 1 << iota
	File
)

type Node struct {
	PackageName string
	Path        string
	Parent      *Node
	Children    []*Node
	Imports     []Edge
	Type        Type
}

func NewNode(packageName string, path string, parent *Node, _type Type) *Node {
	if strings.HasPrefix(path, "./") {
		path = path[2:]
	}

	return &Node{
		PackageName: packageName,
		Path:        path,
		Parent:      parent,
		Type:        _type,
	}
}

func (n *Node) AddImport(imported *Node) {
	edge := NewEdge(n, imported)
	n.Imports = append(n.Imports, edge)

	node := n
	for node.Parent != nil {
		node = node.Parent
		edge := NewEdge(n, imported)
		node.Imports = append(node.Imports, edge)
	}
}

func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
}

type Edge struct {
	From *Node
	To   *Node
}

func NewEdge(from *Node, to *Node) Edge {
	return Edge{From: from, To: to}
}
